package main

import (
	"context"
	ds "miniKV/grpc_gen/dataService"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// rpc连接 将地址替换为leader节点
	conn, err := grpc.Dial("127.0.0.1:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Error(err)
		return
	}
	d := ds.NewDataServiceClient(conn)
	// set
	resp1, err := d.Set(context.Background(), &ds.SetReq{
		Key:   "testKey2",
		Value: "testVal2",
		Info: &ds.ReqInfo{
			ClientId: "1001",
			SeqId:    "101",
		},
	})
	logrus.Infof("set: %+v err: %+v", resp1, err)
	// get exist
	resp2, err := d.Get(context.Background(), &ds.GetReq{
		Key: "testKey2",
	})
	logrus.Infof("get exist: %+v err: %+v", resp2, err)
	// get not exist
	resp3, err := d.Get(context.Background(), &ds.GetReq{
		Key: "notExist",
	})
	logrus.Infof("get not exist: %+v err: %+v", resp3, err)
	// del
	resp4, err := d.Del(context.Background(), &ds.DelReq{
		Key: "testKey2",
		Info: &ds.ReqInfo{
			ClientId: "1001",
			SeqId:    "102",
		},
	})
	logrus.Infof("del: %+v err: %+v", resp4, err)
	// check del
	resp5, err := d.Get(context.Background(), &ds.GetReq{
		Key: "testKey2",
	})
	logrus.Infof("check del: %+v err: %+v", resp5, err)
}
