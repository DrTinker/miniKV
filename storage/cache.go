package storage

type CacheMap struct {
	// key - offset
	keyMap map[string]int64
}

func NewCacheMap() *CacheMap {
	return &CacheMap{
		keyMap: make(map[string]int64),
	}
}

func (c *CacheMap) Put(key string, offset int64) error {
	c.keyMap[key] = offset
	return nil
}

func (c *CacheMap) Del(key string) error {
	delete(c.keyMap, key)
	return nil
}

func (c *CacheMap) Get(key string) int64 {
	return c.keyMap[key]
}

func (c *CacheMap) Exist(key string) bool {
	_, ok := c.keyMap[key]
	return ok
}
