package service

import (
	"miniKV/algo/raft"
	ds "miniKV/grpc_gen/dataService"
	"miniKV/models"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type DataServer struct {
	dataHandler *DataHandler
	dataServer  *grpc.Server
	dataListen  *net.Listener
}

func StartDataServer(me, maxraftstate int, node *raft.RaftNode, applyCh chan models.ApplyMsg,
	snapshotCh chan models.SnapshotMsg) *DataServer {
	s := &DataServer{}
	// 先监听端口，如果失败则没有后续过程
	dataListen, err := net.Listen("tcp", "127.0.0.1:8002")
	if err != nil {
		logrus.Errorf("[KVServer] failed to listen: %v", err)
		return nil
	}
	s.dataListen = &dataListen
	// 创建grpc server
	s.dataHandler = NewDataHandler(me, maxraftstate, node, applyCh, snapshotCh)
	s.dataServer = grpc.NewServer()
	ds.RegisterDataServiceServer(s.dataServer, s.dataHandler)
	// 绑定监听端口
	s.dataServer.Serve(dataListen)

	return s
}
