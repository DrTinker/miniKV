package helper

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
)

// 打开本地文件
// OpenFile 判断文件是否存在  存在则OpenFile 不存在则Create
func OpenFile(path, name string) (*os.File, error) {
	filename := filepath.Join(path, name)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// 先创建目录
		os.MkdirAll(path, os.ModeDir)
		return os.Create(filename) //创建文件
	}
	return os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
}

func ReadFile(path, name string) ([]byte, error) {
	f, err := OpenFile(path, name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func WriteFile(path, name string, data []byte) error {
	filename := filepath.Join(path, name)
	// 目录不存在则创建目录
	_, err := OpenFile(path, name)
	if err != nil {
		return err
	}
	err = os.WriteFile(filename, data, 0666) //写入文件(字节数组)
	if err != nil {
		return err
	}
	return nil
}

// 删除文件
func DelFile(path string) error {
	err := os.Remove(path)
	return err
}

// 返回目录下全部文件列表，按照数字大小排序
func ReadDir(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	res := make([]string, len(list))
	for _, f := range list {
		n, _ := strconv.Atoi(f.Name())
		// 分片从1开始
		res[n-1] = f.Name()
	}
	return res, nil
}

// 文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { //文件或者目录存在
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RemoveDir(path string) error {
	exist, err := PathExists(path)
	if err != nil {
		return err
	}
	if exist {
		err = os.RemoveAll(path)
		if err != nil {
			return err
		}
	}
	return nil
}
