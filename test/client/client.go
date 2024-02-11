package main

import "miniKV/client"

func main() {
	cli := client.NewKVClient("127.0.0.1:8000")
	cli.Connect()
	defer cli.Close()
	cli.Run()
}
