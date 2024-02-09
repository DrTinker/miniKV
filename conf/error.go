package conf

import (
	"errors"
)

// raft
var RaftLogDecodeErr = errors.New("raft log decode error")

var ServerInternalErr = errors.New("ServerInternalErr")
var KeyNotExistErr = errors.New("KeyNotExistErr")
var WrongLeaderErr = errors.New("WrongLeaderErr")
var TimeoutErr = errors.New("TimeoutErr")
