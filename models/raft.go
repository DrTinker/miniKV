package models

import rs "miniKV/grpc_gen/raftService"

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For PartD:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

type LogEntry struct {
	Term         int
	CommandValid bool
	Command      string
}

func (l *LogEntry) ToRPC() *rs.LogEntry {
	res := &rs.LogEntry{
		Term:         int64(l.Term),
		CommandValid: l.CommandValid,
		Command:      l.Command,
	}

	return res
}

func (l *LogEntry) FromRPC(r *rs.LogEntry) {
	l.Term = int(r.Term)
	l.Command = r.Command
	l.CommandValid = r.CommandValid
}

// 投票
type RequestVoteArgs struct {
	Term        int
	CandidateId int
	// 用于竞选安全的参数，任期相同时应通过此参数判断谁的日志更新
	LastLogTerm  int
	LastLogIndex int
}

func (r *RequestVoteArgs) ToRPC() *rs.RequestVoteReq {
	res := &rs.RequestVoteReq{
		Term:         int64(r.Term),
		CandidateId:  int64(r.CandidateId),
		LastLogTerm:  int64(r.LastLogTerm),
		LastLogIndex: int64(r.LastLogIndex),
	}

	return res
}

func (r *RequestVoteArgs) FromRPC(req *rs.RequestVoteReq) {
	r.Term = int(req.Term)
	r.CandidateId = int(req.CandidateId)
	r.LastLogIndex = int(req.LastLogIndex)
	r.LastLogTerm = int(req.LastLogTerm)
}

type RequestVoteReply struct {
	Term int
	// 是否投票
	VoteGranted bool
}

func (r *RequestVoteReply) ToRPC() *rs.RequestVoteResp {
	res := &rs.RequestVoteResp{
		Term:        int64(r.Term),
		VoteGranted: r.VoteGranted,
	}

	return res
}

func (r *RequestVoteReply) FromRPC(resp *rs.RequestVoteResp) {
	r.Term = int(resp.Term)
	r.VoteGranted = resp.VoteGranted
}

// 日志
type AppendEntriesArgs struct {
	Term     int
	LeaderId int

	PrevLogIndex int
	PrevLogTerm  int
	Entries      []LogEntry

	LeaderCommit int
}

func (a *AppendEntriesArgs) ToRPC() *rs.AppendEntriesReq {
	Entries := make([]*rs.LogEntry, len(a.Entries))
	for i, e := range a.Entries {
		Entries[i] = e.ToRPC()
	}
	res := &rs.AppendEntriesReq{
		Term:         int64(a.Term),
		LeaderId:     int64(a.LeaderId),
		PrevLogIndex: int64(a.PrevLogIndex),
		PrevLogTerm:  int64(a.PrevLogTerm),
		LeaderCommit: int64(a.LeaderCommit),
		Entries:      Entries,
	}

	return res
}

func (a *AppendEntriesArgs) FromRPC(r *rs.AppendEntriesReq) {
	a.Term = int(r.Term)
	a.LeaderId = int(r.LeaderId)
	a.PrevLogIndex = int(r.PrevLogIndex)
	a.PrevLogTerm = int(r.PrevLogTerm)
	a.LeaderCommit = int(r.LeaderCommit)

	Entries := make([]LogEntry, len(r.Entries))
	for i, e := range r.Entries {
		Entries[i] = LogEntry{}
		Entries[i].FromRPC(e)
	}
}

type AppendEntriesReply struct {
	Term    int
	Success bool

	// 用于日志和leader不匹配时快速回退
	ConfilictIndex int
	ConfilictTerm  int
}

func (a *AppendEntriesReply) ToRPC() *rs.AppendEntriesResp {
	res := &rs.AppendEntriesResp{
		Term:           int64(a.Term),
		Success:        a.Success,
		ConfilictIndex: int64(a.ConfilictIndex),
		ConfilictTerm:  int64(a.ConfilictTerm),
	}

	return res
}

func (a *AppendEntriesReply) FromRPC(r *rs.AppendEntriesResp) {
	a.Term = int(r.Term)
	a.Success = r.Success
	a.ConfilictIndex = int(r.ConfilictIndex)
	a.ConfilictTerm = int(r.ConfilictTerm)
}

// 快照
type InstallSnapshotArgs struct {
	Term     int
	LeaderId int

	// Snapshot中最后一条日志的index和Term
	LastIncludedIndex int
	LastIncludedTerm  int

	Snapshot []byte
}

func (i *InstallSnapshotArgs) ToRPC() *rs.InstallSnapshotReq {
	res := &rs.InstallSnapshotReq{
		Term:              int64(i.Term),
		LeaderId:          int64(i.LeaderId),
		LastIncludedIndex: int64(i.LastIncludedIndex),
		LastIncludedTerm:  int64(i.LastIncludedTerm),
		Snapshot:          i.Snapshot,
	}

	return res
}

func (i *InstallSnapshotArgs) FromRPC(r *rs.InstallSnapshotReq) {
	i.Term = int(r.Term)
	i.LeaderId = int(r.LeaderId)
	i.LastIncludedIndex = int(r.LastIncludedIndex)
	i.LastIncludedTerm = int(r.LastIncludedTerm)
	i.Snapshot = r.Snapshot
}

type InstallSnapshotReply struct {
	Term int
}

func (i *InstallSnapshotReply) ToRPC() *rs.InstallSnapshotResp {
	res := &rs.InstallSnapshotResp{
		Term: int64(i.Term),
	}

	return res
}

func (i *InstallSnapshotReply) FromRPC(r *rs.InstallSnapshotResp) {
	i.Term = int(r.Term)
}
