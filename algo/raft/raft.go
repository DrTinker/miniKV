package raft

import (
	"miniKV/conf"
	rs "miniKV/grpc_gen/raftService"
	"miniKV/helper"
	"miniKV/models"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO 实现raft算法
// raft模块特有的models
// 节点角色
type Role string

const (
	Follower  Role = "Follower"
	Candidate Role = "Candidate"
	Leader    Role = "Leader"
)

type RaftNode struct {
	mu        sync.Mutex             // 整个节点的锁
	peers     []rs.RaftServiceClient // 集群其他节点的RPC地址
	persister *persister             // Object to hold this peer's persisted state
	me        int                    // this peer's index into peers[]
	dead      int32                  // set by Kill()

	// 节点当前角色
	role Role
	// term 和 vote 需要持久化
	currentTerm int
	votedFor    int

	// 选举相关
	electionStart   time.Time
	electionTimeout time.Duration

	// 日志需要持久化
	log RaftLog

	// 用于日志同步的两个view
	nextIndex  []int
	matchIndex []int

	// apply相关
	commitIndex int                  // 当前节点最后一条提交的log index
	lastApplied int                  // 最后一条应用到状态机的log index
	applyCond   *sync.Cond           // 提交后唤醒apply loop持久化日志
	applyCh     chan models.ApplyMsg // 向应用层传递提交日志的channel

	// 复用applyCh的标志
	snapPending bool
}

func (rf *RaftNode) GetState() (int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	return rf.currentTerm, rf.role == Leader
}

// 是否超过法定人数
func (r *RaftNode) isQuorum(n int) bool {
	return n >= ((len(r.peers) + 1) / 2)
}

// 关闭节点
func (rf *RaftNode) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
}

// 检查节点是否运行
func (rf *RaftNode) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// 检查角色和term
func (rf *RaftNode) contextLostLocked(role Role, term int) bool {
	return !(rf.currentTerm == term && rf.role == role)
}

// 获取match中多数匹配的日志index
func (rf *RaftNode) getMajorityIndexLocked() int {
	// deep copy match后排序
	tmp := make([]int, len(rf.matchIndex))
	copy(tmp, rf.matchIndex)
	sort.Ints(sort.IntSlice(tmp))
	majorityIdx := (len(tmp) - 1) / 2
	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Match index after sort: %v, majority[%d]=%d"),
		tmp, majorityIdx, tmp[majorityIdx])
	return tmp[majorityIdx] // min -> max
}

// 实现状态迁移的三个函数
func (rf *RaftNode) becomeFollowerLocked(term int) {
	if term < rf.currentTerm {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm,
			"Can't become Follower, lower term"))
		return
	}

	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "%s -> Follower, For T%d->T%d"),
		rf.role, rf.currentTerm, term)

	// important! Could only reset the `votedFor` when term increased
	shouldPersist := term != rf.currentTerm
	if term > rf.currentTerm {
		rf.votedFor = -1
	}
	rf.role = Follower
	rf.currentTerm = term
	// 如果任期变化则需要持久化
	if shouldPersist {
		rf.persistLocked()
	}
}

func (rf *RaftNode) becomeCandidateLocked() {
	if rf.role == Leader {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "Leader can't become Candidate"))
		return
	}

	logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "%s -> Candidate, For T%d->T%d"),
		rf.role, rf.currentTerm, rf.currentTerm+1)
	rf.role = Candidate
	rf.currentTerm++
	rf.votedFor = rf.me
	// term和voteFor变化，需要持久化
	rf.persistLocked()
}

func (rf *RaftNode) becomeLeaderLocked() {
	if rf.role != Candidate {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm,
			"%s, Only candidate can become Leader"), rf.role)
		return
	}

	logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "%s -> Leader, For T%d"),
		rf.role, rf.currentTerm)
	rf.role = Leader

	// 初始化next和match两个view
	// next初始设为leader本地log的下一个，如果不匹配会在同步日志时递归的回退
	// match初始为0，代表没有匹配
	for i := 0; i < len(rf.peers); i++ {
		rf.nextIndex[i] = int(rf.log.size() - 1)
		rf.matchIndex[i] = 0
	}
}

// 应用层通过start方法来追加log(执行command)
func (rf *RaftNode) Start(command string) (int, int, bool) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	// 检查是否为leader
	if rf.role != Leader {
		return 0, 0, false
	}
	rf.log.appendFrom(-1, []models.LogEntry{
		{
			CommandValid: true,
			Command:      command,
			Term:         rf.currentTerm,
		},
	})
	// 日志变更，触发持久化
	rf.persistLocked()
	logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "Leader accept log [%d]T%d"),
		rf.log.size()-1, rf.currentTerm)

	return rf.log.size() - 1, rf.currentTerm, true
}

// 创算法层节点
func InitRaftNode(peers []rs.RaftServiceClient, me int, applyCh chan models.ApplyMsg) *RaftNode {
	rf := &RaftNode{}
	rf.peers = peers
	rf.persister = &persister{}
	rf.me = me

	// Your initialization code here (PartA, PartB, PartC).
	rf.log = *NewLog(conf.InvalidIndex, conf.InvalidTerm, nil, nil)

	rf.matchIndex = make([]int, len(rf.peers))
	rf.nextIndex = make([]int, len(rf.peers))

	rf.applyCh = applyCh
	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.applyCond = sync.NewCond(&rf.mu)

	rf.snapPending = false

	// 这里用来从磁盘读取持久化的状态，来恢复raft节点
	rf.readPersist(rf.persister.ReadRaftState())

	// start ticker goroutine to start elections
	go rf.electionTicker()
	go rf.applicationTicker()

	return rf
}
