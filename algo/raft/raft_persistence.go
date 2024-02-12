package raft

import (
	"bytes"
	"miniKV/helper"

	"github.com/sirupsen/logrus"
)

// 根据原论文，需要持久化的有term，voteFor，log

// raft持久化相关函数
func (rf *RaftNode) persistLocked() {
	buf := new(bytes.Buffer)

	e := helper.NewEncoder(buf)
	// encode term vote
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	// encode log除了snapshot
	rf.log.encodeLog(e)

	// 持久化raftstate和snapshot
	raftstate := buf.Bytes()
	rf.persister.Save(raftstate, rf.log.snapshot)
}

// restore previously persisted state.
func (rf *RaftNode) readPersist(data []byte) {
	if data == nil || len(data) < 1 {
		return
	}

	var currentTerm int
	var votedFor int

	buf := bytes.NewBuffer(data)
	d := helper.NewDecoder(buf)
	// 读取term
	if err := d.Decode(&currentTerm); err != nil {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "Read currentTerm error: %v"), err)
		return
	}
	rf.currentTerm = currentTerm

	// 读取vote
	if err := d.Decode(&votedFor); err != nil {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "Read votedFor error: %v"), err)
		return
	}
	rf.votedFor = votedFor

	// 读取log
	if err := rf.log.decodeLog(d); err != nil {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "Read log error: %v"), err)
		return
	}
	// 读取snapshot
	rf.log.snapshot = rf.persister.ReadSnapshot()

	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Read Persist %v"), rf.role)
}
