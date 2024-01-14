package config

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

//
const (
	PUT uint16 = iota
	DEL
)
