package models

// MyParser对应解析出的命令
// Data最长为3 SET K1 V1
type MyCmd struct {
	Data []string
	Err  error
}
