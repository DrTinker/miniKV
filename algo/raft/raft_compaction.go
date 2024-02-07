package raft

import (
	"context"
	"miniKV/helper"
	"miniKV/models"

	"github.com/sirupsen/logrus"
)

// 应用层调用此接口意味着应用层已经完成了一次snapshot，index和snapshot都是应用层传来的参数
func (rf *RaftNode) Snapshot(index int, snapshot []byte) {
	// Your code here (PartD).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Snap on %d"), index)
	// 不能对已经snapshot过的日志再次操作 也 不能对未提交的日志snapshot
	if index <= rf.log.snapLastIdx || index > rf.commitIndex {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Could not snapshot beyond [%d, %d]"), rf.log.snapLastIdx+1, rf.commitIndex)
		return
	}

	rf.log.doSnapshot(index, snapshot)
	rf.persistLocked()
}

func (rf *RaftNode) sendInstallSnapshot(server int, args *models.InstallSnapshotArgs, reply *models.InstallSnapshotReply) bool {
	resp, err := rf.peers[server].InstallSnapshot(context.Background(), args.ToRPC())
	if err != nil {
		logrus.Errorf(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, Lost or crashed"), server)
		return false
	}
	reply.FromRPC(resp)
	return true
}

// leader向其他peer发送snapshot
func (rf *RaftNode) installOnPeer(peer, term int, args *models.InstallSnapshotArgs) {
	reply := &models.InstallSnapshotReply{}
	ok := rf.sendInstallSnapshot(peer, args, reply)
	if !ok {
		return
	}

	rf.mu.Lock()
	defer rf.mu.Unlock()

	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "-> S%d, InstallSnap, Reply=%v"), peer, reply)

	// align the term
	if reply.Term > rf.currentTerm {
		rf.becomeFollowerLocked(reply.Term)
		return
	}
	// 如果有多个 InstallSnapshotReply 乱序回来，且
	// 较小的 args.LastIncludedIndex 后回来的话，如果不加判断，
	// 会造成matchIndex 和 nextIndex 的反复横跳。
	if args.LastIncludedIndex > rf.matchIndex[peer] { // to avoid disorder reply
		// 要更新next和match，否则将不断重复发 InstallSnapshot RPC
		rf.matchIndex[peer] = args.LastIncludedIndex
		rf.nextIndex[peer] = args.LastIncludedIndex + 1
	}

	// note: we need not try to update the commitIndex again,
	// because the snapshot included indexes are all committed
}

// peer接受snapshot后的处理
func (rf *RaftNode) InstallSnapshot(args *models.InstallSnapshotArgs, reply *models.InstallSnapshotReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, RecvSnap, Args=%v"),
		args.LeaderId, args)

	reply.Term = rf.currentTerm
	// align the term
	if args.Term < rf.currentTerm {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, Reject Snap, Higher Term, T%d>T%d"),
			args.LeaderId, rf.currentTerm, args.Term)
		return
	}
	if args.Term > rf.currentTerm {
		rf.becomeFollowerLocked(args.Term)
	}

	// check if it is a RPC which is out of order
	if rf.log.snapLastIdx >= args.LastIncludedIndex {
		logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "<- S%d, Reject Snap, Already installed, Last: %d>=%d"),
			args.LeaderId, rf.log.snapLastIdx, args.LastIncludedIndex)
		return
	}
	// install the snapshot
	// 在算法层的raft结构体中删除对应部分日志
	rf.log.installSnapshot(args.LastIncludedIndex, args.LastIncludedTerm, args.Snapshot)
	// 持久化算法层数据
	rf.persistLocked()
	// 设置标志，避免应用snapshot和应用日志冲突
	rf.snapPending = true
	// 唤醒applyCh进行消息写入，通知应用层
	rf.applyCond.Signal()
}
