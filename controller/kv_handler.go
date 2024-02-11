package controller

import (
	"context"
	"io"
	"miniKV/conf"
	ds "miniKV/grpc_gen/dataService"
	"miniKV/helper"
	"miniKV/parser"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// 实现handler接口，注入tcp server
type KVHandler struct {
	client ds.DataServiceClient
	resCh  chan string
	errCh  chan error
	dead   int32
}

func NewKVHandler(conn *grpc.ClientConn, errCh chan error) KVHandler {
	return KVHandler{
		client: ds.NewDataServiceClient(conn),
		resCh:  make(chan string),
		errCh:  errCh,
	}
}

func (kv *KVHandler) Handle(ctx context.Context, conn net.Conn) {
	parser := parser.NewMyParser(conn)
	cmdCh := parser.GetCmdChan()
	// defer parser.CloseCmdChan()
	go parser.ParseStream()
	// for range 自动判断chan是否关闭
	for {
		select {
		case myCmd := <-cmdCh:
			// 错误处理
			err := myCmd.Err
			if err != nil {
				if err == io.EOF {
					logrus.Error("[KVHandler] connection close")
					return
				}
				logrus.Errorf("[KVHandler] parse cmd err: %v", err)
				continue
			}
			// cmd处理
			cmd := myCmd.Data
			logrus.Infof("[KVHandler] get cmd: %v", cmd)
			go kv.invoke(cmd)
			select {
			case res := <-kv.resCh:
				// 写回结果
				conn.Write([]byte(res))
			case <-time.After(conf.ClientRequestTimeout):
				conn.Write([]byte("request time out"))
			}
		default:
			if kv.isClosed() {
				return
			}
		}
	}
}

func (kv *KVHandler) Close() {
	atomic.StoreInt32(&kv.dead, 1)
}

func (kv *KVHandler) isClosed() bool {
	z := atomic.LoadInt32(&kv.dead)
	return z == 1
}

func (kv *KVHandler) invoke(cmd []string) {
	// 心跳处理
	if len(cmd) == 1 && cmd[0] == conf.HeartBeatArg {
		kv.resCh <- conf.HeartBeatReply
		return
	}
	// 长度不对
	if len(cmd) > 3 || len(cmd) < 2 {
		kv.resCh <- conf.CmdInvaildErr.Error()
		return
	}
	opStr := strings.ToUpper(cmd[0])
	// 无效操作符
	op, ok := conf.OpTable[opStr]
	if !ok {
		kv.resCh <- conf.CmdInvaildErr.Error()
		return
	}
	clientId := helper.GenClientId()
	seqId := helper.GenSeqId()
	// rpc调用
	switch op {
	case conf.OpSet:
		if len(cmd) != 3 {
			kv.resCh <- conf.CmdInvaildErr.Error()
			return
		}
		_, err := kv.client.Set(context.Background(),
			&ds.SetReq{Key: cmd[1], Value: cmd[2],
				Info: &ds.ReqInfo{ClientId: clientId, SeqId: seqId}})
		// leader变更则想外层通知
		if err == conf.WrongLeaderErr {
			kv.errCh <- err
		}
		// 其他错误直接返回
		if err != nil {
			kv.resCh <- err.Error()
			return
		}
		kv.resCh <- "success"
	case conf.OpGet:
		if len(cmd) != 2 {
			kv.resCh <- conf.CmdInvaildErr.Error()
			return
		}
		resp, err := kv.client.Get(context.Background(),
			&ds.GetReq{Key: cmd[1]})
		if err == conf.WrongLeaderErr {
			kv.errCh <- err
			return
		}
		if err != nil {
			kv.resCh <- err.Error()
			return
		}
		kv.resCh <- resp.Value
	case conf.OpDel:
		if len(cmd) != 2 {
			kv.resCh <- conf.CmdInvaildErr.Error()
			return
		}
		resp, err := kv.client.Del(context.Background(),
			&ds.DelReq{Key: cmd[1], Info: &ds.ReqInfo{ClientId: clientId, SeqId: seqId}})
		if err == conf.WrongLeaderErr {
			kv.errCh <- err
			return
		}
		if err != nil {
			kv.resCh <- err.Error()
			return
		}
		kv.resCh <- resp.Value
	}
}
