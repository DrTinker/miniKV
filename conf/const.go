package conf

import (
	"os"
	"path/filepath"
	"time"
)

// entry头大小 = keySize(uint32) + valueSize(uint32) + type(uint16)
const EntryHeaderSize = uint32(10)

// 数据文件默认名称
const DataFileName = "minikv_data"

// merge生成临时文件
const MergeTmpFileName = "minikv_tmp"

// 默认路径
const DiskDefaultPath = "data"

// 大小
const KB = 1024
const MB = 1024 * 1024
const GB = 1024 * 1024 * 1024

// storage 命令
const (
	PUT uint16 = iota
	DEL
)

// tcp 服务器状态
const (
	// 开启
	OPEN uint32 = iota
	// 完全关闭，禁止新链接，清除旧链接
	CLOSE
	// 挂起，禁止新链接，旧链接保留
	HUP
)

const DefaultConnTimeout = 10 * time.Second

// raft
// 实现随机时间间隔选举
const (
	ElectionTimeoutMin time.Duration = 250 * time.Millisecond
	ElectionTimeoutMax time.Duration = 400 * time.Millisecond

	// 心跳发送间隔一定要小于超时时间
	ReplicateInterval time.Duration = 70 * time.Millisecond
)

// 日志不匹配时快速回退
const (
	InvalidIndex int = 0
	InvalidTerm  int = 0
)

// 应用层
// raftstate 和 snapshot的存储地址
const RaftPersistPath = "raft_persistence"
const ServicePersistPath = "service_persistence"
const RaftStateFileName = "state.info"
const RaftSnapFileName = "snapshot.info"

const StateMachineName = "raft"

// 操作类型
type OpType uint8

const (
	OpGet OpType = iota
	OpSet
	OpDel
)

var OpTable = map[string]OpType{
	"SET": OpSet,
	"GET": OpGet,
	"DEL": OpDel,
}

// 操作超时时间
const ClientRequestTimeout = 500 * time.Millisecond

// client关键字
var History_fn = filepath.Join(os.TempDir(), ".liner_example_history")
var KeyWords = []string{"set", "get", "del"}

const HeartBeatArg = "ping"
const HeartBeatReply = "pong"

const CLRF = "\r\n"
