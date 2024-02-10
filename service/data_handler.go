package service

import (
	"bytes"
	"context"
	"fmt"
	"miniKV/algo/raft"
	"miniKV/conf"
	ds "miniKV/grpc_gen/dataService"
	"miniKV/helper"
	"miniKV/models"
	"miniKV/storage"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type DataHandler struct {
	mu   sync.Mutex
	me   int
	node *raft.RaftNode
	// 与算法层交互
	applyCh    chan models.ApplyMsg
	snapshotCh chan models.SnapshotMsg
	// snapshot阈值
	maxraftstate int
	// 算法层回调的通知channel
	callbackChs map[int]chan models.CallBackMsg
	// 确保请求幂等性
	duplicateTable map[int]models.LastOpInfo
	// 记录已经应用到状态机的最后一条index
	lastApplied int
	// 状态机
	stateMachine *storage.DB
}

func NewDataHandler(me, maxraftstate int, node *raft.RaftNode, applyCh chan models.ApplyMsg,
	snapshotCh chan models.SnapshotMsg) *DataHandler {
	d := &DataHandler{}
	d.me = me
	d.maxraftstate = maxraftstate
	d.node = node
	d.applyCh = applyCh
	d.snapshotCh = snapshotCh
	d.callbackChs = make(map[int]chan models.CallBackMsg)
	d.lastApplied = 0
	d.duplicateTable = make(map[int]models.LastOpInfo)

	d.restoreFromSnapshot(nil)

	go d.getRaftReply()
	return d
}

func (d *DataHandler) Get(ctx context.Context, req *ds.GetReq) (resp *ds.GetResp, err error) {
	resp = &ds.GetResp{}
	op := models.Op{
		Key:    req.Key,
		OpType: conf.OpGet,
	}
	command, err1 := op.Encode()
	if err1 != nil {
		logrus.Errorf("[DataHandler] encode command err: %v", err1)
		return resp, conf.ServerInternalErr
	}
	// 交付raft算法层
	index, _, isLeader := d.node.Start(command)
	// 强制主读，解决raft同步延时问题
	if !isLeader {
		logrus.Errorf("[DataHandler] command on wrong leader: %v", op)
		return resp, conf.WrongLeaderErr
	}
	// 获取callback
	d.mu.Lock()
	cch := d.getCallbackCh(index)
	d.mu.Unlock()

	// 等待结果
	select {
	case res := <-cch:
		if res.Err != nil {
			logrus.Errorf("[DataHandler] execute command %v err: %v", op, res.Err)
			err = res.Err
		} else {
			resp.Success = true
			resp.Value = res.Value
		}
	case <-time.After(conf.ClientRequestTimeout):
		logrus.Errorf("[DataHandler] execute command %v timeout", op)
		err = conf.TimeoutErr
	}

	// 销毁ch
	go func() {
		d.mu.Lock()
		d.removeCallbackCh(index)
		d.mu.Unlock()
	}()

	return resp, err
}

func (d *DataHandler) Set(ctx context.Context, req *ds.SetReq) (resp *ds.SetResp, err error) {
	resp = &ds.SetResp{}
	clientId, seqId := int(req.Info.ClientId), int(req.Info.SeqId)
	// 判断幂等性
	d.mu.Lock()
	if d.checkIdempotent(clientId, seqId) {
		lastRes := d.duplicateTable[clientId]
		// 重复直接返回上次存储的结果
		resp.Success = lastRes.Reply.Success
		d.mu.Unlock()
		return resp, lastRes.Reply.Err
	}
	d.mu.Unlock()

	op := models.Op{
		Key:      req.Key,
		Value:    req.Value,
		OpType:   conf.OpSet,
		ClientId: clientId,
		SeqId:    seqId,
	}
	command, err1 := op.Encode()
	if err1 != nil {
		logrus.Errorf("[DataHandler] encode command err: %v", err1)
		return resp, conf.ServerInternalErr
	}
	// 交付raft算法层
	index, _, isLeader := d.node.Start(command)
	// 写请求必须在leader上执行
	if !isLeader {
		logrus.Errorf("[DataHandler] command on wrong leader: %v", op)
		return resp, conf.WrongLeaderErr
	}
	// 获取callback
	d.mu.Lock()
	cch := d.getCallbackCh(index)
	d.mu.Unlock()

	// 等待结果
	select {
	case res := <-cch:
		if res.Err != nil {
			logrus.Errorf("[DataHandler] execute command %v err: %v", op, res.Err)
			err = conf.ServerInternalErr
		} else {
			resp.Success = true
		}
	case <-time.After(conf.ClientRequestTimeout):
		logrus.Errorf("[DataHandler] execute command %v timeout", op)
		err = conf.TimeoutErr
	}

	// 销毁ch
	go func() {
		d.mu.Lock()
		d.removeCallbackCh(index)
		d.mu.Unlock()
	}()

	return resp, err
}

func (d *DataHandler) Del(ctx context.Context, req *ds.DelReq) (resp *ds.DelResp, err error) {
	resp = &ds.DelResp{}
	clientId, seqId := int(req.Info.ClientId), int(req.Info.SeqId)
	// 判断幂等性
	d.mu.Lock()
	if d.checkIdempotent(clientId, seqId) {
		lastRes := d.duplicateTable[clientId]
		// 重复直接返回上次存储的结果
		resp.Success = lastRes.Reply.Success
		resp.Value = lastRes.Reply.Value
		d.mu.Unlock()
		return resp, lastRes.Reply.Err
	}
	d.mu.Unlock()
	// 封装operation
	op := models.Op{
		Key:      req.Key,
		OpType:   conf.OpDel,
		ClientId: clientId,
		SeqId:    seqId,
	}
	command, err1 := op.Encode()
	if err1 != nil {
		logrus.Errorf("[DataHandler] encode command err: %v", err1)
		return resp, conf.ServerInternalErr
	}
	// 交付raft算法层
	index, _, isLeader := d.node.Start(command)
	// 写请求必须在leader上执行
	if !isLeader {
		logrus.Errorf("[DataHandler] command on wrong leader: %v", op)
		return resp, conf.WrongLeaderErr
	}
	// 获取callback
	d.mu.Lock()
	cch := d.getCallbackCh(index)
	d.mu.Unlock()

	// 等待结果
	select {
	case res := <-cch:
		if res.Err != nil {
			logrus.Errorf("[DataHandler] execute command %v err: %v", op, res.Err)
			err = conf.ServerInternalErr
		} else {
			resp.Success = true
			resp.Value = res.Value
		}
	case <-time.After(conf.ClientRequestTimeout):
		logrus.Errorf("[DataHandler] execute command %v timeout", op)
		err = conf.TimeoutErr
	}

	// 销毁ch
	go func() {
		d.mu.Lock()
		d.removeCallbackCh(index)
		d.mu.Unlock()
	}()

	return resp, err
}

// 监听channel，获取算法层的执行结果
func (d *DataHandler) getRaftReply() {
	for {
		select {
		case msg := <-d.applyCh:
			d.handleApplyMsg(&msg)
		case msg := <-d.snapshotCh:
			d.handleSnapshotMsg(&msg)
		}
	}
}

func (d *DataHandler) handleApplyMsg(msg *models.ApplyMsg) {
	// 命令无效直接返回
	if !msg.CommandValid {
		return
	}
	// 判断是否处理过
	d.mu.Lock()
	if msg.CommandIndex <= d.lastApplied {
		d.mu.Unlock()
		return
	}
	d.lastApplied = msg.CommandIndex
	// 应用
	// decode operation
	cmdData := msg.Command
	op := &models.Op{}
	op.Decode(cmdData)
	// 应用状态机
	reply := d.applyToStateMachine(op)
	// 保存本次请求
	if op.OpType != conf.OpGet {
		d.duplicateTable[op.ClientId] = models.LastOpInfo{
			SeqId: op.SeqId,
			Reply: reply,
		}
	}
	// 结果发送callback
	if _, isLeader := d.node.GetState(); isLeader {
		notifyCh := d.getCallbackCh(msg.CommandIndex)
		notifyCh <- reply
	}
	// 判断是否需要snapshot
	if d.node.GetRaftStateSize() >= d.maxraftstate {
		d.makeSnapshot(msg.CommandIndex)
	}
	d.mu.Unlock()
}

func (d *DataHandler) handleSnapshotMsg(msg *models.SnapshotMsg) {
	d.mu.Lock()
	d.restoreFromSnapshot(msg.Snapshot)
	d.lastApplied = msg.SnapshotIndex
	d.mu.Unlock()
}

func (d *DataHandler) applyToStateMachine(op *models.Op) models.CallBackMsg {
	var val string
	var err error
	switch op.OpType {
	case conf.OpGet:
		val, err = d.stateMachine.Get(op.Key)
	case conf.OpSet:
		err = d.stateMachine.Put(op.Key, op.Value)
	case conf.OpDel:
		val, err = d.stateMachine.Del(op.Key)
	}
	msg := models.CallBackMsg{}
	if err != nil && err != conf.KeyNotExistErr {
		// 一般错误
		logrus.Errorf("[DataHandler] applyToStateMachine err: %v", err)
		msg.Err = conf.ServerInternalErr
		msg.Success = false
	} else {
		// 成功或者key不存在
		msg.Success = true
		msg.Err = err
		msg.Value = val
	}
	logrus.Infof("[DataHandler] applyToStateMachine info: %v", msg)
	return msg
}

// 获取callback channel
func (d *DataHandler) getCallbackCh(index int) chan models.CallBackMsg {
	if _, ok := d.callbackChs[index]; !ok {
		d.callbackChs[index] = make(chan models.CallBackMsg, 1)
	}
	return d.callbackChs[index]
}

// 销毁callback channel
func (d *DataHandler) removeCallbackCh(index int) {
	delete(d.callbackChs, index)
}

// 检查请求幂等性
func (d *DataHandler) checkIdempotent(clientId, seqId int) bool {
	res, ok := d.duplicateTable[clientId]
	return ok && seqId <= res.SeqId
}

func (d *DataHandler) makeSnapshot(index int) {
	buf := new(bytes.Buffer)
	// 保存
	enc := helper.NewEncoder(buf)
	_ = enc.Encode(d.duplicateTable)
	err := helper.WriteFile(conf.ServicePersistPath, conf.RaftSnapFileName, buf.Bytes())
	if err != nil {
		logrus.Errorf("[DataHandler] makeSnapshot db err: %+v", err)
		return
	}
	// merge lsm tree
	err = d.stateMachine.Merge()
	if err != nil {
		logrus.Errorf("[DataHandler] makeSnapshot db err: %+v", err)
		return
	}
	d.node.Snapshot(index, buf.Bytes())
}

func (d *DataHandler) restoreFromSnapshot(snapshot []byte) {
	// 恢复幂等表
	dupTable := make(map[int]models.LastOpInfo)
	var buf *bytes.Buffer
	// 为空则从本地磁盘磁盘读取
	if len(snapshot) == 0 {
		data, err := helper.ReadFile(conf.ServicePersistPath, conf.RaftSnapFileName)
		if err != nil {
			panic(fmt.Sprintf("failed to restore state from snapshpt, err: %v", err))
		}
		snapshot = data
	}
	if len(snapshot) != 0 {
		buf = bytes.NewBuffer(snapshot)
		dec := helper.NewDecoder(buf)
		if dec.Decode(&dupTable) != nil {
			panic("failed to restore state from snapshpt")
		}
	}

	// 恢复db
	stateMachine, err := storage.OpenDB(conf.StateMachineName)
	if err != nil {
		panic(fmt.Sprintf("failed to restore state from snapshpt, err: %v", err))
	}

	d.stateMachine = stateMachine
	d.duplicateTable = dupTable
}
