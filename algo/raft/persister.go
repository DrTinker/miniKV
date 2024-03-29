package raft

import (
	"miniKV/conf"
	"miniKV/helper"
	"sync"

	"github.com/sirupsen/logrus"
)

// raft持久化工具类
type persister struct {
	mu            sync.Mutex
	raftStateSize int
}

func MakePersister() *persister {
	return &persister{}
}

func (ps *persister) ReadRaftState() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	data, err := helper.ReadFile(conf.RaftPersistPath, conf.RaftStateFileName)
	if err != nil {
		logrus.Errorf("[Raft] ReadRaftState err: %v", err)
		return nil
	}
	return data
}

// 保存state和snapshot地址
func (ps *persister) Save(raftstate []byte, snapshot []byte) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	err := helper.WriteFile(conf.RaftPersistPath, conf.RaftStateFileName, raftstate)
	if err != nil {
		logrus.Errorf("[Raft] Save err: %v", err)
	}
	err = helper.WriteFile(conf.RaftPersistPath, conf.RaftSnapFileName, snapshot)
	if err != nil {
		logrus.Errorf("[Raft] Save err: %v", err)
	}
	// 记录state大小
	ps.raftStateSize = len(raftstate)
}

func (ps *persister) ReadSnapshot() []byte {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	data, err := helper.ReadFile(conf.RaftPersistPath, conf.RaftSnapFileName)
	if err != nil {
		logrus.Errorf("[Raft] ReadSnapshot err: %v", err)
		return nil
	}
	return data
}
