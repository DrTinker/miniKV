package raft

import (
	"context"
	"miniKV/conf"
	"miniKV/helper"
	"miniKV/models"
	"time"

	"github.com/sirupsen/logrus"
)

func (rf *RaftNode) sendAppendEntries(server int, args *models.AppendEntriesArgs, reply *models.AppendEntriesReply) bool {
	resp, err := rf.peers[server].AppendEntries(context.Background(), args.ToRPC())
	if err != nil {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, Lost or crashed"), server)
		return false
	}
	reply.FromRPC(resp)
	return true
}

// 返回是否继续当前term的
func (rf *RaftNode) startReplilcation(term int) bool {
	replicateToPeer := func(peer int, args *models.AppendEntriesArgs) {
		reply := &models.AppendEntriesReply{}
		ok := rf.sendAppendEntries(peer, args, reply)
		if !ok {
			return
		}

		rf.mu.Lock()
		defer rf.mu.Unlock()

		// 收到更大任期退回Follower
		if reply.Term > rf.currentTerm {
			rf.becomeFollowerLocked(reply.Term)
			return
		}

		// check context lost
		if rf.contextLostLocked(Leader, term) {
			logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, Context Lost, T%d:Leader->T%d:%s"), peer, term, rf.currentTerm, rf.role)
			return
		}

		// 处理reply
		if !reply.Success {
			prevNext := rf.nextIndex[peer]
			// 如果ConfilictTerm为空，说明follower日志差的太多，直接无视term回退到
			// follower的最后一条日志
			if reply.ConfilictTerm == conf.InvalidTerm {
				rf.nextIndex[peer] = reply.ConfilictIndex
			} else {
				// 如果不为空，说明term不匹配，直接回退到冲突的term的第一条日志
				firstTermIndex := rf.log.firstLogFor(reply.ConfilictTerm)
				if firstTermIndex != conf.InvalidIndex {
					rf.nextIndex[peer] = firstTermIndex + 1
				} else {
					rf.nextIndex[peer] = reply.ConfilictIndex
				}
			}
			// 设置在回退日志匹配的过程中，nextIndex只能单调递减
			// 避免多个reply乱序到达，使得nextIndex反复横跳
			if rf.nextIndex[peer] > prevNext {
				rf.nextIndex[peer] = prevNext
			}
			logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Log not matched in %d, Update next=%d"),
				args.PrevLogIndex, rf.nextIndex[peer])
			return
		}
		// 日志同步成功，更新两个view
		rf.matchIndex[peer] = args.PrevLogIndex + len(args.Entries)
		rf.nextIndex[peer] = rf.matchIndex[peer] + 1

		// 更新commitIndex
		// If there exists an N such that N>commitIndex, a majority of matchIndex[i]≥N,
		// and log[N].term==currentTerm: set commitIndex=N
		n := rf.getMajorityIndexLocked()
		if n > rf.commitIndex && rf.log.at(n).Term == rf.currentTerm {
			logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Leader update the commit index %d->%d"),
				rf.commitIndex, n)
			rf.commitIndex = n
			rf.applyCond.Signal()
		}
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if rf.contextLostLocked(Leader, term) {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Leader[T%d] -> %s[T%d]"),
			term, rf.role, rf.currentTerm)
		return false
	}

	for peer := 0; peer < len(rf.peers); peer++ {
		if peer == rf.me {
			rf.matchIndex[peer] = rf.log.size() - 1
			rf.nextIndex[peer] = rf.log.size()
			continue
		}

		// 计算prev
		prevLogIndex := rf.matchIndex[peer]
		// 判断prev是否被snapshot截断，如果截断需要发送installSnap RPC
		// prevLogIndex = rf.log.snapLastIdx会取到dummy节点，也是可用的
		if prevLogIndex < rf.log.snapLastIdx {
			args := &models.InstallSnapshotArgs{
				Term:              rf.currentTerm,
				LeaderId:          rf.me,
				LastIncludedIndex: rf.log.snapLastIdx,
				LastIncludedTerm:  rf.log.snapLastTerm,
				Snapshot:          rf.log.snapshot,
			}
			logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, InstallSnap, Args=%v"),
				peer, args)
			go rf.installOnPeer(peer, term, args)
			continue
		}
		prevLogTerm := rf.log.at(prevLogIndex).Term
		args := &models.AppendEntriesArgs{
			Term:         term,
			LeaderId:     rf.me,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entries:      rf.log.tail(prevLogIndex + 1),
			LeaderCommit: rf.commitIndex,
		}

		go replicateToPeer(peer, args)
	}

	return true
}

// 生命周期是一个term，term变了这个ticker就结束
func (rf *RaftNode) replicationTicker(term int) {
	for !rf.killed() {
		// 如果超出当前任期或者角色发生变化，则不再维护心跳
		if ok := rf.startReplilcation(term); !ok {
			return
		}

		time.Sleep(conf.ReplicateInterval)
	}
}

func (rf *RaftNode) AppendEntries(args *models.AppendEntriesArgs, reply *models.AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, Receive log, Prev=[%d]T%d, Len()=%d"),
		args.LeaderId, args.PrevLogIndex, args.PrevLogTerm, len(args.Entries))
	// replay initialized
	reply.Term = rf.currentTerm
	reply.Success = false

	// term对齐
	// Leader的term小于自己，直接拒绝
	if args.Term < rf.currentTerm {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, Reject log"), args.LeaderId)
		return
	}

	// 无论是否接收日志，只要term对齐，就需要reset
	defer rf.resetElectionTimerLocked()

	// 大于自己，变为follower
	if args.Term >= rf.currentTerm {
		rf.becomeFollowerLocked(args.Term)
	}

	// 接收日志，根据PrevLogIndex判断匹配点
	// 不匹配，先判断index是否越界(说明本节点日志比leader少)，不少的话再判断term是否相同
	if args.PrevLogIndex >= rf.log.size() {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, Reject Log, Follower log too short, Len:%d <= Prev:%d"),
			args.LeaderId, rf.log.size(), args.PrevLogIndex)
		// 日志差的较多，
		reply.ConfilictIndex = rf.log.size()
		reply.ConfilictTerm = conf.InvalidTerm
		return
	}
	//  同一index的日志term也要相同
	if args.PrevLogTerm != rf.log.at(args.PrevLogIndex).Term {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, Reject Log, Prev log not match, [%d]: T%d != T%d"),
			args.LeaderId, args.PrevLogIndex, rf.log.at(args.PrevLogIndex).Term, args.PrevLogTerm)
		// term不匹配时，应通知leader回退一个term的日志
		reply.ConfilictTerm = rf.log.at(args.PrevLogIndex).Term
		reply.ConfilictIndex = rf.log.firstLogFor(reply.ConfilictTerm)
		return
	}
	// 下来是匹配的情况，复制日志到本节点
	// 日志覆盖，体现了raft的strong leader特点，与主节点不一致的日志全背覆盖
	rf.log.appendFrom(args.PrevLogIndex, args.Entries)
	reply.Success = true
	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Follower append logs: (%d, %d]"),
		args.PrevLogIndex, args.PrevLogIndex+len(args.Entries))
	// log变化，需要持久化
	rf.persistLocked()
	// 日志应用
	if args.LeaderCommit > rf.commitIndex {
		// commitIndex变为LeaderCommit和args中日志条目最高索引 二者中的较小值
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Follower update the commit index %d->%d"),
			rf.commitIndex, args.LeaderCommit)
		rf.commitIndex = args.LeaderCommit
		// 这里args中的日志已经被复制到rf.log中了，所以只要比较 args中日志条目最高索引 = len(rf.log)
		if rf.commitIndex >= rf.log.size() {
			rf.commitIndex = rf.log.size() - 1
		}
		rf.applyCond.Signal()
	}
}
