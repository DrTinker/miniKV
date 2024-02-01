package models

import (
	"miniKV/helper"
	"net"
)

// 对tcp链接的封装，使其可以等待一定时间再被释放
type TcpClient struct {
	Conn net.Conn
	Wait helper.Wait
}
