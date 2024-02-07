package raft

import (
	"miniKV/helper"
	"miniKV/models"

	"github.com/sirupsen/logrus"
)

// 实现apply loop
func (rf *RaftNode) applicationTicker() {
	for !rf.killed() {
		rf.mu.Lock()
		// cond的作用
		// cond在初始化时需要注入一把锁，在这里时rf.mu
		// 调用wait时，会将当前goroutine挂到对应事件的等待队列中，
		// 同时解锁rf.mu并将当前goroutine阻塞在wait这一行
		// 当外部goroutine调用signal时，会唤醒当前goroutine，并且重新获得rf.mu
		// 因此这里不满足条件时不会占有锁，同时因为唤醒时自动加锁，所以下面要有解锁语句
		rf.applyCond.Wait()
		// 封装要应用的日志
		entries := make([]models.LogEntry, 0)
		snapPendingInstall := rf.snapPending
		if !snapPendingInstall {
			// leader通知snapshot的rpc先来，此时部分还未apply的日志可能被snapshot
			// 再收到日志同步请求，要应用日志时，就会出现lastApplied < log.snapLastIdx
			// 因为leader已经snapshot的日志必然被leaderapply了，因此follower也应视作apply了这部分日志
			// lastApplied = log.snapLastIdx 的意思就是apply了[lastApplied, log.snapLastIdx]部分的日志
			if rf.lastApplied < rf.log.snapLastIdx {
				rf.lastApplied = rf.log.snapLastIdx
			}

			// 处理[log.snapLastIdx, commitIndex]部分的日志apply
			start := rf.lastApplied + 1
			end := rf.commitIndex
			if end >= rf.log.size() {
				end = rf.log.size() - 1
			}
			for i := start; i <= end; i++ {
				entries = append(entries, rf.log.at(i))
			}
		}
		rf.mu.Unlock()

		// apply在非snapPendingInstall得情况进行日志的应用
		// 而snapPendingInstall=true时，需要先进行snapshot的应用
		if !snapPendingInstall {
			for i, entry := range entries {
				rf.applyCh <- models.ApplyMsg{
					CommandValid: entry.CommandValid,
					Command:      entry.Command,
					CommandIndex: rf.lastApplied + 1 + i, // must be cautious
				}
			}
		} else {
			rf.applyCh <- models.ApplyMsg{
				SnapshotValid: true,
				Snapshot:      rf.log.snapshot,
				SnapshotIndex: rf.log.snapLastIdx,
				SnapshotTerm:  rf.log.snapLastTerm,
			}
		}

		rf.mu.Lock()
		if !snapPendingInstall {
			// 日志应用
			logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Apply log for [%d, %d]"),
				rf.lastApplied+1, rf.lastApplied+len(entries))
			rf.lastApplied += len(entries)
		} else {
			// sanpshot intall
			logrus.Infof(helper.RaftPrefix(rf.me, rf.currentTerm, "Install Snapshot for [0, %d]"),
				rf.log.snapLastIdx)
			rf.lastApplied = rf.log.snapLastIdx
			if rf.commitIndex < rf.lastApplied {
				rf.commitIndex = rf.lastApplied
			}
			rf.snapPending = false
		}
		rf.mu.Unlock()
	}
}
