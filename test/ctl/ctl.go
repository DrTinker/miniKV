package main

import "miniKV/controller"

func main() {
	configs := []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}
	ctl := controller.NewKVController("127.0.0.1", 8000, configs)
	ctl.Serve()
}
