package raft

import (
	"errors"
	"fmt"
	"miniKV/conf"
	"miniKV/helper"
	"miniKV/models"
)

type RaftLog struct {
	// snapshot最后一条日志的idx和term
	snapLastIdx  int
	snapLastTerm int

	// snapshot的一个磁盘地址 [1, snapLastIdx] 磁盘io太耗时
	// 还是在内存中保存
	snapshot []byte
	// 内存中保存的日志，初始有一个空的dummy节点
	tailLog []models.LogEntry
}

func NewLog(snapLastIdx, snapLastTerm int, snapshot []byte, entries []models.LogEntry) *RaftLog {
	rl := &RaftLog{
		snapLastIdx:  snapLastIdx,
		snapLastTerm: snapLastTerm,
		snapshot:     snapshot,
	}

	// +1时因为有一个头结点
	rl.tailLog = make([]models.LogEntry, 0, 1+len(entries))
	rl.tailLog = append(rl.tailLog, models.LogEntry{
		Term: snapLastTerm,
	})
	rl.tailLog = append(rl.tailLog, entries...)

	return rl
}

// decode raft log
func (rl *RaftLog) decodeLog(d *helper.MyDecoder) error {
	var lastIdx int
	if err := d.Decode(&lastIdx); err != nil {
		return errors.New("[Raft] log decode last include index failed")
	}
	rl.snapLastIdx = lastIdx

	var lastTerm int
	if err := d.Decode(&lastTerm); err != nil {
		return errors.New("[Raft] log decode last include term failed")
	}
	rl.snapLastTerm = lastTerm

	var snapshotAddr string
	if err := d.Decode(&snapshotAddr); err != nil {
		return errors.New("[Raft] log decode last include term failed")
	}

	var log []models.LogEntry
	if err := d.Decode(&log); err != nil {
		return errors.New("[Raft] log ecode tail log failed")
	}
	rl.tailLog = log

	return nil
}

// encode raft log
func (rl *RaftLog) encodeLog(e *helper.MyEncoder) {
	e.Encode(rl.snapLastIdx)
	e.Encode(rl.snapLastTerm)
	e.Encode(rl.tailLog)
}

// 以下为封装堆tailLog的操作，主要是简化snapshot截断后的下标处理

// 获取全部日志大小(包含snapshot的部分)
func (rl *RaftLog) size() int {
	return rl.snapLastIdx + len(rl.tailLog)
}

// 逻辑下标转换实际下标
func (rl *RaftLog) idx(logicIdx int) int {
	// 越界
	if logicIdx < rl.snapLastIdx || logicIdx >= int(rl.size()) {
		panic(fmt.Sprintf("%d is out of index [%d, %d]", logicIdx, rl.snapLastIdx+1, rl.size()-1))
	}
	return int(logicIdx - rl.snapLastIdx)
}

// 取下标
func (rl *RaftLog) at(logicIdx int) models.LogEntry {
	return rl.tailLog[rl.idx(logicIdx)]
}

// 获取最后一条日志的idx和term
func (rl *RaftLog) last() (idx, term int) {
	return rl.size() - 1, rl.tailLog[len(rl.tailLog)-1].Term
}

// 截取后半部分
func (rl *RaftLog) tail(startIdx int) []models.LogEntry {
	if startIdx >= rl.size() {
		return nil
	}
	return rl.tailLog[rl.idx(startIdx):]
}

// 寻找当前term第一条log
func (rl *RaftLog) firstLogFor(term int) int {
	for idx, entry := range rl.tailLog {
		if entry.Term == term {
			return int(idx) + rl.snapLastIdx
		} else if entry.Term > term {
			break
		}
	}
	return conf.InvalidIndex
}

// append prevIdx=-1表示直接在后面append
func (rl *RaftLog) appendFrom(prevIdx int, entries []models.LogEntry) {
	if prevIdx == -1 {
		rl.tailLog = append(rl.tailLog, entries...)
		return
	}
	rl.tailLog = append(rl.tailLog[:rl.idx(prevIdx)+1], entries...)
}

// 执行应用层传来的snapshot命令
func (rl *RaftLog) doSnapshot(index int, snapshot []byte) {
	// 逻辑下标转换实际下标
	idx := rl.idx(index)

	rl.snapLastTerm = rl.tailLog[idx].Term
	rl.snapLastIdx = index
	// 记录snapshot的磁盘存储位置
	rl.snapshot = snapshot

	// raft层只需要丢弃index之前的日志即可
	newLog := make([]models.LogEntry, 0, rl.size()-rl.snapLastIdx)
	newLog = append(newLog, models.LogEntry{
		Term: rl.snapLastTerm,
	})
	newLog = append(newLog, rl.tailLog[idx+1:]...)
	rl.tailLog = newLog
}

// 应用snapshot
func (rl *RaftLog) installSnapshot(index, term int, snapshot []byte) {
	rl.snapLastIdx = index
	rl.snapLastTerm = term
	// 存snapshot操作在应用层进行过了，这里只需要记录addr
	rl.snapshot = snapshot

	// Follower 将 snapshot 保存到rf.log 中时，完全覆盖掉以前日志就可以。
	// 因为新来的 snapshot 的最后一条日志下标（ lastIncludeIndex ）一定是大于原来
	// Follower log 的最后一条 index 的（即 Leader 发过来的 snapshot 肯定包含更多信息），
	// 否则 Leader 就不需要发送 snapshot 给 Follower 了。
	newLog := make([]models.LogEntry, 0, 1)
	newLog = append(newLog, models.LogEntry{
		Term: rl.snapLastTerm,
	})
	rl.tailLog = newLog
}
