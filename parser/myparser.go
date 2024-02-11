package parser

import (
	"bufio"
	"io"
	"miniKV/models"
	"strings"
)

// 一个简易的应用层命令解析协议
// 只解析单行数据
type MyParser struct {
	reader io.Reader
	cmdCh  chan *models.MyCmd
}

func NewMyParser(reader io.Reader) *MyParser {
	cmdCh := make(chan *models.MyCmd)
	mp := MyParser{}
	mp.cmdCh = cmdCh
	mp.reader = reader

	return &mp
}

func (mp *MyParser) GetCmdChan() <-chan *models.MyCmd {
	return mp.cmdCh
}

func (mp *MyParser) CloseCmdChan() {
	close(mp.cmdCh)
}

func (mp *MyParser) ParseStream() {
	// 使用 bufio 标准库提供的缓冲区功能
	reader := bufio.NewReader(mp.reader)
	for {
		// ReadString 会一直阻塞直到遇到分隔符 '\n'
		// 遇到分隔符后会返回上次遇到分隔符或连接建立后收到的所有数据, 包括分隔符本身
		// 若在遇到分隔符之前遇到异常, ReadString 会返回已收到的数据和错误信息
		msg, err := reader.ReadString('\n')
		if msg == "" {
			return
		}
		// raw = SET K V\r\n
		// msg 应当为 SET K V 或 DEL K 或 GET K
		cmds := strings.Split(msg[:len(msg)-2], " ")
		// 封装为MyCmd
		mycmd := models.MyCmd{
			Data: cmds,
			Err:  err,
		}
		// 写入channel
		mp.cmdCh <- &mycmd
	}
}
