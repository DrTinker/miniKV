package main

import "miniKV/network"

func main() {
	// 测试echo
	network.ListenAndServe("127.0.0.1:8001")

	// 测试server
	// server := network.NewTcpServer("127.0.0.1", 8000, network.EchoHandler{})
	// server.StartServer()
}
