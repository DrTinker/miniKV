package service

import (
	"miniKV/algo/raft"
	rs "miniKV/grpc_gen/raftService"
	"miniKV/models"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// 实现raft的应用层，为了减小锁的粒度，将应用层拆分为两个server
// raftServer：开放8001端口，处理raft的3个rpc请求
// dataServer：开放8002端口，处理controller发来的command
// 如果锁在外层，那么会同时锁住
type KVService struct {
	me int
	rf *raft.RaftNode
	// raft rpc相关
	peerAddrs  []string
	raftServer *RaftServer
	dataServer *DataServer
	// kvServer   *KVServer
	// 与算法层交互
	applyCh    chan models.ApplyMsg
	snapshotCh chan models.SnapshotMsg
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
	// 建立和其他peer的链接
	peers := make([]rs.RaftServiceClient, len(peerAddrs))
	for i, addr := range peerAddrs {
		conn, err := grpc.Dial(addr)
		if err != nil {
			logrus.Errorf("[KVServer] raft rpc connect %s err: %v", addr, err)
			continue
		}
		peers[i] = rs.NewRaftServiceClient(conn)
	}
	// 创建和算法层交互的channel
	applyCh := make(chan models.ApplyMsg)
	s.applyCh = applyCh
	snapshotCh := make(chan models.SnapshotMsg)
	s.snapshotCh = snapshotCh
	// 创建算法层node
	rf := raft.InitRaftNode(peers, me, applyCh, snapshotCh)
	s.rf = rf
	// 创建raftServer
	s.raftServer = StartRaftServer(rf)
	// 创建dataServer
	s.dataServer = StartDataServer(me, maxraftstate, rf, applyCh, snapshotCh)

	return s
}
