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
	SeqId string
	Reply CallBackMsg
}

type Op struct {
	Key      string
	Value    string
	OpType   conf.OpType
	ClientId string
	SeqId    string
}

func (o *Op) Encode() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Op) Decode(data []byte) {
	json.Unmarshal(data, o)
}
