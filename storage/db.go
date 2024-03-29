package storage

import (
	"errors"
	"io"
	"miniKV/conf"
	"os"
	"path/filepath"
	"sync"
)

type DB struct {
	// 磁盘文件
	DiskFile DiskFile
	// 内存索引
	CacheMap CacheMap
	// 磁盘存储路径
	DirPath string
	// 文件排他所
	lock sync.RWMutex
}

func initDB(name string) (*DB, error) {
	// 创建内存map
	cache := NewCacheMap()

	dir := filepath.Join(conf.DiskDefaultPath, name)
	// 如果数据库目录不存在，则新建一个
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 创建磁盘文件
	disk, err := NewDiskFile(dir)
	if err != nil {
		return nil, err
	}

	return &DB{
		DiskFile: *disk,
		CacheMap: *cache,
		DirPath:  dir,
	}, nil
}

func NewDB(name string) (*DB, error) {
	return initDB(name)
}

func OpenDB(name string) (*DB, error) {
	db, err := initDB(name)
	if err != nil {
		return nil, err
	}
	// 加载cache
	err = db.loadCacheMap()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d *DB) Put(key, value string) error {
	// 加锁
	d.lock.Lock()
	defer d.lock.Unlock()
	// 包装entry
	keyData, valData := []byte(key), []byte(value)
	e := NewEntry(keyData, valData, conf.PUT)
	// 记录插入前的offset
	offset := d.DiskFile.GetOffset()
	// 先更新磁盘
	err := d.DiskFile.Write(e)
	if err != nil {
		return err
	}
	// 再更新内存
	d.CacheMap.Put(key, offset)
	return nil
}

// TODO 返回删除的value
func (d *DB) Del(key string) (string, error) {
	// 判断是否存在
	ok := d.CacheMap.Exist(key)
	if !ok {
		return "", conf.KeyNotExistErr
	}

	// 加锁
	d.lock.Lock()
	defer d.lock.Unlock()

	// 先读取原数据
	offset := d.CacheMap.Get(key)
	// 根据offset从内存获取entry
	ori, err := d.DiskFile.Read(offset)
	if err != nil && err != io.EOF {
		return "", err
	}

	// 在磁盘写入del entry
	e := NewEntry([]byte(key), nil, conf.DEL)
	err = d.DiskFile.Write(e)
	if err != nil {
		return "", err
	}

	// 只在内存删除
	d.CacheMap.Del(key)

	if ori == nil {
		return "", nil
	}
	return string(ori.Value), nil
}

// val值，key是否存在，err
func (d *DB) Get(key string) (string, error) {
	// 先查询内存
	ok := d.CacheMap.Exist(key)
	if !ok {
		return "", conf.KeyNotExistErr
	}

	// 从内存获取offset
	offset := d.CacheMap.Get(key)

	// 加锁
	d.lock.Lock()
	defer d.lock.Unlock()

	// 根据offset从内存获取entry
	e, err := d.DiskFile.Read(offset)
	if err != nil && err != io.EOF {
		return "", err
	}
	if e != nil {
		return string(e.Value), nil
	}

	return "", nil
}

func (d *DB) Merge() error {
	vaildEntries := []*Entry{}
	offset := int64(0)

	// 加锁
	d.lock.Lock()
	defer d.lock.Unlock()

	// 从原文件中筛选出有效的entry
	for {
		// 从磁盘读取
		e, err := d.DiskFile.Read(int64(offset))
		if err != nil {
			// 已经遍历全部entry
			if err == io.EOF {
				break
			}
			// 遍历时出错
			return err
		}
		// 去内存中查询，内存中保存着最新的key
		// 存在这个key 而且内存中记录的offset和磁盘中的一致，说明当前entry有效
		if d.CacheMap.Exist(string(e.Key)) && d.CacheMap.Get(string(e.Key)) == offset {
			vaildEntries = append(vaildEntries, e)
		}
		offset += e.GetSize()
	}

	if len(vaildEntries) == 0 {
		return errors.New("empty disk file")
	}

	// 创建临时文件
	mFile, err := NewMergeDiskFile(d.DirPath)
	if err != nil {
		return err
	}
	// 临时文件一定要删除
	defer os.Remove(mFile.GetFileName())

	// 将有效entry全部写入新文件
	for _, ve := range vaildEntries {
		writeOff := mFile.GetOffset()
		// 写入
		err := mFile.Write(ve)
		if err != nil {
			return err
		}
		// 更新内存
		d.CacheMap.Put(string(ve.Key), writeOff)
	}

	// 用临时文件代替原文件
	oriFileName, tmpFileName := d.DiskFile.GetFileName(), mFile.GetFileName()
	// 删除旧的数据文件
	d.DiskFile.Close()
	err = os.Remove(oriFileName)
	if err != nil {
		return err
	}
	// 临时文件变更为新的数据文件
	mFile.Close()
	err = os.Rename(tmpFileName, d.DirPath+string(os.PathSeparator)+conf.DataFileName)
	if err != nil {
		return err
	}
	d.DiskFile = *mFile

	return nil
}

func (d *DB) Close() error {
	return d.DiskFile.Close()
}

func (d *DB) loadCacheMap() error {
	pos := int64(0)
	// 读取全部entry，记录offset
	for pos < d.DiskFile.GetOffset() {
		// 读取entry
		e, err := d.DiskFile.Read(pos)
		if err != nil {
			return err
		}
		// 写入cache
		d.CacheMap.Put(string(e.Key), pos)
		// 获取entry长度，移动指针
		pos += e.GetSize()
	}

	return nil
}
