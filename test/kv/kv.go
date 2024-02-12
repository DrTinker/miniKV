package main

import (
	"miniKV/conf"
	"miniKV/service"
	"miniKV/start"
	"sync"
)

func main() {
	start.InitLog()
	// new三个server出来
	peerAddrs := []string{"127.0.0.1:8001", "127.0.0.1:8002", "127.0.0.1:8003"}
	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		s := service.NewKVServer(peerAddrs, i, conf.MB)
		go func(s *service.KVService) {
			go s.ConnectToPeers()
			s.Serve()
			wg.Done()
		}(s)
	}
	wg.Wait()
}
