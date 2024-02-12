package service

import (
	"context"
	"miniKV/algo/raft"
	rs "miniKV/grpc_gen/raftService"
	"miniKV/models"
)

// 实现raftService
type RaftHandler struct {
	node *raft.RaftNode
}

func NewRaftHandler(node *raft.RaftNode) *RaftHandler {
	proxy := &RaftHandler{}
	proxy.node = node
	return proxy
}

func (r *RaftHandler) RequestVote(ctx context.Context, req *rs.RequestVoteReq) (resp *rs.RequestVoteResp, err error) {
	resp = &rs.RequestVoteResp{}
	args := &models.RequestVoteArgs{}
	args.FromRPC(req)

	reply := &models.RequestVoteReply{}
	reply.FromRPC(resp)

	r.node.RequestVote(args, reply)
	resp = reply.ToRPC()
	return resp, nil
}

func (r *RaftHandler) AppendEntries(ctx context.Context, req *rs.AppendEntriesReq) (resp *rs.AppendEntriesResp, err error) {
	resp = &rs.AppendEntriesResp{}
	args := &models.AppendEntriesArgs{}
	args.FromRPC(req)

	reply := &models.AppendEntriesReply{}
	reply.FromRPC(resp)

	r.node.AppendEntries(args, reply)
	resp = reply.ToRPC()
	return resp, nil
}

func (r *RaftHandler) InstallSnapshot(ctx context.Context, req *rs.InstallSnapshotReq) (resp *rs.InstallSnapshotResp, err error) {
	resp = &rs.InstallSnapshotResp{}
	args := &models.InstallSnapshotArgs{}
	args.FromRPC(req)

	reply := &models.InstallSnapshotReply{}
	reply.FromRPC(resp)

	r.node.InstallSnapshot(args, reply)
	resp = reply.ToRPC()
	return resp, nil
}
