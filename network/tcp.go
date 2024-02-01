package network

import (
	"context"
	"fmt"
	"miniKV/conf"
	"miniKV/interface/tcp"
	"miniKV/models"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// 实现一个tcp服务器
type TcpServer struct {
	// ip
	ip string
	// 开放端口
	port int
	// listener
	listener net.Listener
	// 关闭channel，监听关闭时间
	closeChan chan struct{}
	// 状态标识
	state atomic.Uint32
	// 记录连接的map
	activeConn sync.Map
	// handler
	handler tcp.Handler
}

func NewTcpServer(ip string, port int, handler tcp.Handler) *TcpServer {
	// 初始化map
	server := &TcpServer{
		ip:      ip,
		port:    port,
		handler: handler,
	}

	return server
}

func (t *TcpServer) StartServer() {
	address := fmt.Sprintf("%s:%d", t.ip, t.port)
	// 初始化链接
	l, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Errorf("Start tcp server err: %v", err)
	}
	logrus.Infof("Start tcp server on: %s", address)
	// 初始化closeChan
	cc := make(chan struct{})
	// 设置
	t.listener = l
	t.closeChan = cc
	t.state.Store(conf.OPEN)
	// 监听系统调用
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			t.closeChan <- struct{}{}
		}
	}()
	// 启动
	t.listenAndServe()
}

// 开启对port的监听，并监听closeChan
func (t *TcpServer) listenAndServe() {
	// 监听关闭通知
	go func() {
		<-t.closeChan
		logrus.Info("tcp server shutting down...")
		// 停止监听，listener.Accept()会立即返回 io.EOF
		// 关闭失败重试
		for err := t.Close(); err != nil; {
			logrus.Warnf("tcp server shut down err: %v", err)
			time.Sleep(1 * time.Second)
		}
	}()

	// 在异常退出后释放资源
	defer t.Close()

	ctx := context.Background()
	// 处理链接
	var wg sync.WaitGroup
	for {
		// Accept 会一直阻塞直到有新的连接建立或者listen中断才会返回
		conn, err := t.listener.Accept()
		if err != nil {
			// 通常是由于listener被关闭无法继续监听导致的错误
			logrus.Errorf("accept err: %v", err)
			break
		}
		// 新来的链接
		logrus.Infof("new tcp link from: %v", conn.RemoteAddr())
		wg.Add(1)
		// 开启新的 goroutine 处理该连接
		go func() {
			defer wg.Done()
			t.handler.Handle(ctx, conn)
		}()
	}
	// 阻塞直到链接处理完毕
	wg.Wait()
}

func (t *TcpServer) Close() error {
	// 设置关闭状态
	t.state.Store(conf.CLOSE)
	// 关闭listener
	err := t.listener.Close()
	// 关闭全部链接
	t.activeConn.Range(func(key, value any) bool {
		cli := value.(*models.TcpClient)
		// 等待默认时间关闭
		cli.Wait.WaitWithTimeout(conf.DefaultConnTimeout)
		cli.Conn.Close()
		return true
	})

	return err
}
