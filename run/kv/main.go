package main

import (
	"flag"
	"miniKV/conf"
	"miniKV/controller"
	"miniKV/network"
	"miniKV/service"
	"miniKV/start"
	"strings"
	"sync"
)

const defaultPeers = "127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003"
const defaultCtl = "127.0.0.1:8000"

func init() {
	start.InitLog()
}

func main() {
	// 从参数获取addrs
	peersStr := flag.String("peers", defaultPeers, "set work nodes in cluster split by space, format: 127.0.0.1:8001 127.0.0.1:8002")
	ctlAddr := flag.String("ctl", defaultCtl, "set controller node, format: 127.0.0.1:8000")
	peerAddrs := strings.Split(*peersStr, " ")

	// 启动service各个节点
	n := len(peerAddrs)
	wg := sync.WaitGroup{}
	wg.Add(n + 1)
	for i := 0; i < n; i++ {
		s := service.NewKVServer(peerAddrs, i, conf.MB)
		go func(s *service.KVService) {
			go s.ConnectToPeers()
			s.Serve()
			wg.Done()
		}(s)
	}
	// 启动controller
	go func() {
		ctl := controller.NewKVController(peerAddrs)
		server := network.NewTcpServer(*ctlAddr, ctl)
		server.StartServer()
		wg.Done()
	}()
	wg.Wait()

}
