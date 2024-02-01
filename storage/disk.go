package storage

import (
	"errors"
	"miniKV/conf"
	"os"
	"path/filepath"
	"sync"
)

type DiskFile struct {
	// 文件
	file *os.File
	// 当前写入指针位置
	offset int64
	// header buffer池
	headerBufferPool *sync.Pool
}

func newStruct(dir string) (*DiskFile, error) {
	// 打开磁盘文件，不存在则创建
	file, err := os.OpenFile(dir, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	// 获取文件大小
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// 根据文件大小设置offset
	offset := info.Size()
	// 创建buffer pool
	bufPool := &sync.Pool{
		New: func() interface{} {
			return make([]byte, conf.EntryHeaderSize)
		},
	}

	return &DiskFile{
		file:             file,
		offset:           offset,
		headerBufferPool: bufPool,
	}, nil
}

func NewDiskFile(path string) (*DiskFile, error) {
	dir := filepath.Join(path, conf.DataFileName)
	return newStruct(dir)
}

func NewMergeDiskFile(path string) (*DiskFile, error) {
	dir := filepath.Join(path, conf.MergeTmpFileName)
	return newStruct(dir)
}

func (d *DiskFile) GetOffset() int64 {
	return d.offset
}

func (d *DiskFile) GetFileName() string {
	return d.file.Name()
}

func (d *DiskFile) Close() error {
	if d.file == nil {
		return errors.New("db not exist")
	}

	return d.file.Close()
}

// 根据offset读取
func (d *DiskFile) Read(offset int64) (*Entry, error) {
	// 从pool中取出buffer备用
	buffer := d.headerBufferPool.Get().([]byte)
	// defer放回pool中
	defer d.headerBufferPool.Put(buffer)

	// 先读取header
	n, err := d.file.ReadAt(buffer, offset)
	if err != nil || n != int(conf.EntryHeaderSize) {
		return nil, err
	}
	// 解码header
	entry, err := Decode(buffer)
	if err != nil {
		return nil, err
	}
	// 读取成功offset移动
	offset += int64(n)

	// 读取key
	if entry.KeySize > 0 {
		key := make([]byte, entry.KeySize)
		// 读取key
		_, err = d.file.ReadAt(key, offset)
		if err != nil {
			return nil, err
		}
		entry.Key = key
	}
	offset += int64(entry.KeySize)

	// 读取value
	if entry.ValueSize > 0 {
		value := make([]byte, entry.ValueSize)
		// 读取key
		_, err = d.file.ReadAt(value, offset)
		if err != nil {
			return nil, err
		}
		entry.Value = value
	}

	return entry, nil
}

// 写入
func (d *DiskFile) Write(entry *Entry) error {
	// 编码
	data, err := entry.Encode()
	if err != nil {
		return err
	}
	// 写入
	_, err = d.file.WriteAt(data, d.offset)
	if err != nil {
		return err
	}
	// 修改offset
	d.offset += entry.GetSize()

	return nil
}
