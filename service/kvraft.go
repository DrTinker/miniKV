package service

import (
	"miniKV/algo/raft"
	ds "miniKV/grpc_gen/dataService"
	rs "miniKV/grpc_gen/raftService"
	"miniKV/models"
	"net"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 实现raft的应用层，为了减小锁的粒度，将应用层拆分为两个server
// raftServer：开放8001端口，处理raft的3个rpc请求
// dataServer：开放8002端口，处理controller发来的command
// 如果锁在外层，那么会同时锁住
type KVService struct {
	me int
	rf *raft.RaftNode
	// raft rpc相关
	peerAddrs   []string
	peerClients []*grpc.ClientConn
	server      *grpc.Server
	listen      net.Listener
	// handler
	raftHandler *RaftHandler
	dataHandler *DataHandler
	// kvServer   *KVServer
	// 与算法层交互
	applyCh    chan models.ApplyMsg
	snapshotCh chan models.SnapshotMsg
	// 让算法层等待peer链接成功
	readyCh chan struct{}
	dead    int32
	// snapshot阈值
	maxraftstate int
}

// 初始化server
// peers是peer的address TODO改造成后台线程从controller定期拉取
func NewKVServer(peerAddrs []string, me int, maxraftstate int) *KVService {
	s := &KVService{}

	// 设置基础属性
	s.me = me
	s.maxraftstate = maxraftstate
	s.peerAddrs = peerAddrs

	// 创建和算法层交互的channel
	applyCh := make(chan models.ApplyMsg)
	s.applyCh = applyCh
	snapshotCh := make(chan models.SnapshotMsg)
	s.snapshotCh = snapshotCh

	s.readyCh = make(chan struct{})

	// 建立和其他peer的链接
	peers := make([]*grpc.ClientConn, len(peerAddrs))
	s.peerClients = peers

	// 创建算法层node
	rf := raft.InitRaftNode(s.peerClients, s.me, s.readyCh, s.applyCh, s.snapshotCh)
	s.rf = rf

	// 设置监听端口
	listen, err := net.Listen("tcp", s.peerAddrs[s.me])
	if err != nil {
		logrus.Errorf("[KVServer] failed to listen: %v", err)
		return nil
	}
	s.listen = listen
	s.server = grpc.NewServer()

	// 初始化handler
	s.raftHandler = NewRaftHandler(s.rf)
	s.dataHandler = NewDataHandler(s.me, s.maxraftstate, s.rf, s.applyCh, s.snapshotCh)

	// 注册服务到server
	rs.RegisterRaftServiceServer(s.server, s.raftHandler)
	ds.RegisterDataServiceServer(s.server, s.dataHandler)

	return s
}

// 运行serve方法阻塞main线程
func (kv *KVService) Serve() {
	logrus.Infof("[KVService] serve on %s", kv.peerAddrs[kv.me])
	kv.server.Serve(kv.listen)
}

func (kv *KVService) ConnectToPeers() {
	for !kv.killed() {
		cnt := 0
		// 链接
		for i, addr := range kv.peerAddrs {
			// 初始化
			if kv.peerClients[i] == nil {
				kv.peerClients[i] = &grpc.ClientConn{}
			}
			// 链接
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				logrus.Errorf("[KVServer] raft rpc connect %s err: %v", addr, err)
				continue
			}
			kv.peerClients[i] = conn
			cnt++
		}
		// 超过半数节点连接成功则启动算法层
		if kv.isQuorum(cnt) {
			break
		}
	}
	kv.readyCh <- struct{}{}
}

func (kv *KVService) Kill() {
	atomic.StoreInt32(&kv.dead, 1)
	kv.rf.Kill()

}

func (kv *KVService) killed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == 1
}

// 是否超过法定人数
func (kv *KVService) isQuorum(n int) bool {
	return n >= ((len(kv.peerAddrs) + 1) / 2)
}
