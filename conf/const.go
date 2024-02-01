package conf

import "time"

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
