package models

import (
	"encoding/json"
	"miniKV/conf"
)

type CallBackMsg struct {
	Success bool
	Value   string
	Err     error
}

// 幂等性判断
type LastOpInfo struct {
	SeqId int
	Reply CallBackMsg
}

type Op struct {
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
	Key      string
	Value    string
	OpType   conf.OpType
	ClientId int
	SeqId    int
}

func (o *Op) Encode() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Op) Decode(data []byte) {
	json.Unmarshal(data, o)
}
