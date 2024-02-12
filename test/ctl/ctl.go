package main

import (
	"miniKV/controller"
	"miniKV/network"
)

func main() {
	configs := []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}
	ctl := controller.NewKVController(configs)
	server := network.NewTcpServer("127.0.0.1:8000", ctl)
	server.StartServer()
}
