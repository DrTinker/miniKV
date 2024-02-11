package main

import (
	"flag"
	"miniKV/controller"
	"miniKV/network"
	"strings"
)

const defaultPeers = "127.0.0.1:8001 127.0.0.1:8002 127.0.0.1:8003"
const defaultCtl = "127.0.0.1:8000"

func main() {
	// 从参数获取addrs
	peersStr := flag.String("peers", defaultPeers, "set work nodes in cluster split by space, format: 127.0.0.1:8001 127.0.0.1:8002")
	ctlAddr := flag.String("ctl", defaultCtl, "set controller node, format: 127.0.0.1:8000")
	peerAddrs := strings.Split(*peersStr, " ")

	ctl := controller.NewKVController(peerAddrs)
	server := network.NewTcpServer(*ctlAddr, ctl)
	server.StartServer()
}
