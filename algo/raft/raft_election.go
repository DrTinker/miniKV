package raft

import (
	"context"
	"math/rand"
	"miniKV/conf"
	rs "miniKV/grpc_gen/raftService"
	"miniKV/helper"
	"miniKV/models"
	"time"

	"github.com/sirupsen/logrus"
)

// 重置选举超时计时器
func (rf *RaftNode) resetElectionTimerLocked() {
	rf.electionStart = time.Now()
	randRange := int64(conf.ElectionTimeoutMax - conf.ElectionTimeoutMin)
	rf.electionTimeout = conf.ElectionTimeoutMin + time.Duration(rand.Int63()%randRange)
}

// 判断是否选举超时
func (rf *RaftNode) isElectionTimeoutLocked() bool {
	return time.Since(rf.electionStart) > rf.electionTimeout
}

// 判断谁的日志更新
func (rf *RaftNode) isMoreUpToDateLocked(candidateIndex, candidateTerm int) bool {
	lastIdx, lastTerm := rf.log.last()
	logrus.Debugf(helper.RaftPrefix(rf.me, rf.currentTerm, "Compare last log, Me: [%d]T%d, Candidate: [%d]T%d"),
		lastIdx, lastTerm, candidateIndex, candidateTerm)
	// 比较term，term大的新
	if candidateTerm != lastTerm {
		return lastTerm > candidateTerm
	}
	return lastIdx > candidateIndex
}

func (rf *RaftNode) RequestVote(args *models.RequestVoteArgs, reply *models.RequestVoteReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// align the term
	reply.Term = rf.currentTerm
	if rf.currentTerm > args.Term {
		logrus.Debugf(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, Reject vote, higher term, T%d>T%d"),
			args.CandidateId, rf.currentTerm, args.Term)
		reply.VoteGranted = false
		return
	}
	if rf.currentTerm < args.Term {
		rf.becomeFollowerLocked(args.Term)
	}

	// 检查日志是否更新
	if rf.isMoreUpToDateLocked(args.LastLogIndex, args.LastLogTerm) {
		// 如果本地的更新的话则拒绝投票
		logrus.Debugf(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, Reject Vote, S%d's log less up-to-date"),
			args.CandidateId)
		reply.VoteGranted = false
		return
	}

	// 检查是否投过票
	if rf.votedFor != -1 {
		logrus.Debugf(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, Reject, Already voted S%d"),
			args.CandidateId, rf.votedFor)
		reply.VoteGranted = false
		return
	}

	reply.VoteGranted = true
	rf.votedFor = args.CandidateId
	// voteFor变更需要持久化
	rf.persistLocked()
	// 投票成功才能重置timer
	rf.resetElectionTimerLocked()
	logrus.Debugf(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d"), args.CandidateId)
}

func (rf *RaftNode) sendRequestVote(server int, args *models.RequestVoteArgs, reply *models.RequestVoteReply) bool {
	client := rs.NewRaftServiceClient(rf.peers[server])
	resp, err := client.RequestVote(context.Background(), args.ToRPC())
	if err != nil {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "Ask vote from %v, Lost or error"), server)
		return false
	}
	reply.FromRPC(resp)
	return true
}

func (rf *RaftNode) startElection(term int) bool {
	votes := 0
	askVoteFromPeer := func(peer int, args *models.RequestVoteArgs) {
		reply := &models.RequestVoteReply{}
		// 调用RPC，出错返回
		ok := rf.sendRequestVote(peer, args, reply)
		if !ok {
			return
		}

		// 完成调用之后才能上锁
		rf.mu.Lock()
		defer rf.mu.Unlock()

		// align the term
		if reply.Term > rf.currentTerm {
			rf.becomeFollowerLocked(reply.Term)
			return
		}

		// check the context
		if rf.contextLostLocked(Candidate, term) {
			logrus.Debugf(helper.RaftPrefix(rf.me, rf.currentTerm, "Lost context, abort RequestVoteReply in T%d"), rf.currentTerm)
			return
		}

		// count votes
		if reply.VoteGranted {
			votes++
		}
		if rf.isQuorum(votes) {
			rf.becomeLeaderLocked()
			go rf.replicationTicker(term)
		}
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()

	// every time locked
	if rf.contextLostLocked(Candidate, term) {
		return false
	}

	lastLogIndex, lastLogTerm := rf.log.last()
	for peer := 0; peer < len(rf.peers); peer++ {
		if peer == rf.me {
			votes++
			continue
		}

		args := &models.RequestVoteArgs{
			Term:         term,
			CandidateId:  rf.me,
			LastLogTerm:  lastLogTerm,
			LastLogIndex: lastLogIndex,
		}
		go askVoteFromPeer(peer, args)
	}

	return true
}

func (rf *RaftNode) electionTicker() {
	for !rf.killed() {
		// Your code here (PartA)
		// Check if a leader election should be started.
		rf.mu.Lock()
		if rf.role != Leader && rf.isElectionTimeoutLocked() {
			rf.becomeCandidateLocked()
			go rf.startElection(rf.currentTerm)
		}
		rf.mu.Unlock()

		// pause for a random amount of time between 50 and 350
		// milliseconds.
		ms := 50 + (rand.Int63() % 300)
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}
