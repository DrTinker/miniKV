package network

import (
	"miniKV/helper"
	"net"
)

// server模块不涉及共享的结构体，故将models放在这里

// 对tcp链接的封装，使其可以等待一定时间再被释放
type TcpClient struct {
	Conn net.Conn
	Wait helper.Wait
}
