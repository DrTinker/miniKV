package service

import (
	"miniKV/algo/raft"
	rs "miniKV/grpc_gen/raftService"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RaftServer struct {
	raftProxy  *RaftHandler
	raftServer *grpc.Server
	raftListen *net.Listener
}

func StartRaftServer(node *raft.RaftNode) *RaftServer {
	s := &RaftServer{}
	// 先监听端口，如果失败则没有后续过程
	raftListen, err := net.Listen("tcp", "127.0.0.1:8001")
	if err != nil {
		logrus.Errorf("[KVServer] failed to listen: %v", err)
		return nil
	}
	s.raftListen = &raftListen
	// 创建grpc server
	s.raftProxy = NewRaftHandler(node)
	s.raftServer = grpc.NewServer()
	rs.RegisterRaftServiceServer(s.raftServer, s.raftProxy)
	// 绑定监听端口
	s.raftServer.Serve(raftListen)

	return s
}
