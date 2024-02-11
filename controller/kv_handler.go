package controller

import (
	"context"
	"fmt"
	"io"
	"miniKV/conf"
	ds "miniKV/grpc_gen/dataService"
	"miniKV/helper"
	"miniKV/parser"
	"net"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

// 实现handler接口，注入tcp server
func (kv *KVController) Handle(ctx context.Context, conn net.Conn) {
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
			case err := <-kv.errCh:
				logrus.Info("leader changed")
				// 如果leader发生变更，则获取新leader重试本次cmd
				if err == conf.WrongLeaderErr {
					before := kv.leaderId
					kv.getLeader()
					fmt.Printf("before: %d, after: %d\n", before, kv.leaderId)
				}
				conn.Write([]byte("server busy, please retry"))
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

func (kv *KVController) invoke(cmd []string) {
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
	client := ds.NewDataServiceClient(kv.clusterConn[kv.leaderId])
	switch op {
	case conf.OpSet:
		if len(cmd) != 3 {
			kv.resCh <- conf.CmdInvaildErr.Error()
			return
		}
		_, err := client.Set(context.Background(),
			&ds.SetReq{Key: cmd[1], Value: cmd[2],
				Info: &ds.ReqInfo{ClientId: clientId, SeqId: seqId}})
		// 错误处理
		if err != nil {
			kv.errHandler(err)
			return
		}
		kv.resCh <- "success"
	case conf.OpGet:
		if len(cmd) != 2 {
			kv.resCh <- conf.CmdInvaildErr.Error()
			return
		}
		resp, err := client.Get(context.Background(),
			&ds.GetReq{Key: cmd[1]})
		if err != nil {
			kv.errHandler(err)
			return
		}
		kv.resCh <- resp.Value
	case conf.OpDel:
		if len(cmd) != 2 {
			kv.resCh <- conf.CmdInvaildErr.Error()
			return
		}
		resp, err := client.Del(context.Background(),
			&ds.DelReq{Key: cmd[1], Info: &ds.ReqInfo{ClientId: clientId, SeqId: seqId}})
		if err != nil {
			kv.errHandler(err)
			return
		}
		kv.resCh <- resp.Value
	}
}

// TODO 增加zhuang增加状态码
// grpc 返回的err和原来的err不同
func (kv *KVController) errHandler(err error) {
	fromError, ok := status.FromError(err)
	if !ok {
		logrus.Warnf("[KVController] rpc return a strange err: %v", err)
	}
	switch {
	case fromError.Message() == conf.WrongLeaderErr.Error():
		kv.errCh <- conf.WrongLeaderErr
	case fromError.Message() == conf.KeyNotExistErr.Error():
		kv.resCh <- "key not exist"
	case fromError.Message() == conf.CmdInvaildErr.Error():
		kv.resCh <- "invaild cmd"
	case fromError.Message() == conf.ServerInternalErr.Error():
		kv.resCh <- "server error"
	case fromError.Message() == conf.TimeoutErr.Error():
		logrus.Debug(err)
	}
}
