package main

import (
	"flag"
	"miniKV/client"
)

const defaultCtl = "127.0.0.1:8000"

func main() {
	ctlAddr := flag.String("ctl", defaultCtl, "address of controller node, format: 127.0.0.1:8000")
	cli := client.NewKVClient(*ctlAddr)
	cli.Connect()
	defer cli.Close()
	cli.Run()
}
