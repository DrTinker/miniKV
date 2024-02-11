package client

import (
	"fmt"
	"miniKV/conf"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/peterh/liner"
	"github.com/sirupsen/logrus"
)

// TODO 实现客户端相关逻辑
type KVClient struct {
	remoteAddr string
	conn       net.Conn
	tool       *liner.State
	dead       int32
}

func NewKVClient(addr string) *KVClient {
	cli := &KVClient{}
	cli.remoteAddr = addr
	line := liner.NewLiner()
	// 键盘中断
	line.SetCtrlCAborts(true)
	// 补全逻辑
	line.SetCompleter(func(line string) (c []string) {
		for _, kw := range conf.KeyWords {
			if strings.HasPrefix(kw, strings.ToLower(line)) {
				c = append(c, kw)
			}
		}
		return
	})
	cli.tool = line
	return cli
}

func (kv *KVClient) Connect() {
	// 建立tcp连接
	conn, err := net.Dial("tcp", kv.remoteAddr)
	if err != nil {
		logrus.Errorf("[KVClient] connect to server err: %v", err)
		return
	}
	kv.conn = conn
	atomic.StoreInt32(&kv.dead, int32(conf.OPEN))
}

func (kv *KVClient) Close() {
	// 关闭conn
	kv.conn.Close()
	// 保存命令行历史记录
	if f, err := os.Create(conf.History_fn); err != nil {
		logrus.Infof("Error writing history file: %v", err)
	} else {
		kv.tool.WriteHistory(f)
		f.Close()
	}
	// 关闭命令行工具
	kv.tool.Close()
	// 设置关闭标志
	atomic.StoreInt32(&kv.dead, int32(conf.CLOSE))
}

func (kv *KVClient) Run() {
	// go kv.keepAlive()
	// 运行命令行工具
	for !kv.isClosed() {
		if cmd, err := kv.tool.Prompt(fmt.Sprintf("%s> ", kv.remoteAddr)); err == nil {
			// 记录历史值
			kv.tool.AppendHistory(cmd)

			// 发送cmd
			cmd += conf.CLRF
			_, err1 := kv.conn.Write([]byte(cmd))
			if err1 != nil {
				fmt.Printf("send failed, err: %v\n", err1)
			}

			// 读取返回
			buf := make([]byte, conf.KB)
			_, err := kv.conn.Read(buf)
			if err != nil {
				fmt.Printf("recv failed, err: %v\n", err)
			}
			// 输出返回值
			fmt.Println(string(buf))
		} else if err == liner.ErrPromptAborted {
			fmt.Println("Aborted")
			return
		} else {
			fmt.Printf("Error reading line: %v\n", err)
		}
	}
}

func (kv *KVClient) isClosed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == int32(conf.CLOSE)
}

func (kv *KVClient) keepAlive() {
	for !kv.isClosed() {
		_, err := kv.conn.Write([]byte(conf.HeartBeatArg))
		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			// handle timeout
			kv.Connect()
		}
		time.Sleep(conf.ClientRequestTimeout)
	}
}
