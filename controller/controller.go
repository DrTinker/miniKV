package controller

import (
	"context"
	ds "miniKV/grpc_gen/dataService"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TODO 实现controller，作为集群网关服务
type KVController struct {
	// KV集群内部各个服务addr
	configs     []string
	clusterConn []*grpc.ClientConn
	// 记录leader
	leaderId int
	// 记录hadnler传来的err
	errCh chan error
	resCh chan string
	dead  int32
}

func NewKVController(configs []string) *KVController {
	kv := KVController{}
	// 连接全部节点
	kv.configs = configs
	kv.clusterConn = make([]*grpc.ClientConn, len(configs))
	for i, cfg := range configs {
		// rpc连接 将地址替换为leader节点
		conn, err := grpc.Dial(cfg, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logrus.Errorf("[NewKVController] conn to node %v err: %v", cfg, err)
			continue
		}
		kv.clusterConn[i] = conn
	}
	kv.errCh = make(chan error)
	kv.resCh = make(chan string)
	// 找到leader
	kv.getLeader()

	return &kv
}

func (kv *KVController) Close() {
	atomic.StoreInt32(&kv.dead, 1)
}

func (kv *KVController) isClosed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == 1
}

func (kv *KVController) getLeader() {
	for i, conn := range kv.clusterConn {
		cli := ds.NewDataServiceClient(conn)
		resp, _ := cli.GetState(context.Background(), &ds.GetStateReq{})
		if resp != nil && resp.IsLeader {
			kv.leaderId = i
			return
		}
	}
}

// func (kv *KVController) listenLeaderChange() {
// 	for {
// 		err := <-kv.errCh
// 		if err == conf.WrongLeaderErr {
// 			logrus.Info("wrong leader, try search new leader")
// 			// 获取leader
// 			kv.getLeader()
// 			// 重设handler
// 			handler := NewKVHandler(kv.clusterConn[kv.leaderId], kv.errCh)
// 			kv.server.SetHandler(&handler)
// 		}
// 	}
// }

// 调用者把请求发送到ctl，ctl轮询访问所有config看是不是leader
// 如果是leader则执行这条命令
// 然后记录leaderId

// TODO 让KV集群自动上报leader信息，KV集群开启线程定期拉取所有节点的信息
