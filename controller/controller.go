package controller

import (
	"context"
	"miniKV/network"
	"net"
)

// TODO 实现controller，作为集群网关服务
type KVController struct {
	// 对client暴露的tcp server
	server network.TcpServer
	// KV集群内部各个服务addr
	configs []string
}

// 实现handler接口，注入tcp server
type KVHandler struct {
}

func (kv *KVHandler) Handle(ctx context.Context, conn net.Conn) {
	// 解析tcp包
	// 获取命令
	// 执行对应命令
}

// 调用者把请求发送到ctl，ctl轮询访问所有config看是不是leader
// 如果是leader则执行这条命令
// 然后记录leaderId

// TODO 让KV集群自动上报leader信息，KV集群开启线程定期拉取所有节点的信息
