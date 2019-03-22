package inMemcache

import (
	"errors"
	"fmt"
	"sync"
)

//内存缓存
type inMemcache struct {
	c     map[string][]byte
	mutex sync.RWMutex //枷锁标识
	Stat
}

//实现接口Cacher
func (c *inMemcache) Set(key string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	tmp, exist := c.c[key]
	if exist {
		c.del(key, tmp)
	}

	c.c[key] = v
	c.add(key, v)

	return nil
}

func (c *inMemcache) Get(key string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if v, ok := c.c[key]; ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("current key %s not exist", key))
}

func (c *inMemcache) Delete(key string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if v, ok := c.c[key]; ok {
		delete(c.c, key)
		c.del(key, v)

		return nil
	}

	return errors.New(fmt.Sprintf("current key %s not exist", key))
}

func (c *inMemcache) GetStat() Stat {
	return c.Stat
}

func ToString(v []byte, e error) (string, error) {
	if e != nil {
		return "", e
	}

	return string(v), nil
}

//工厂函数
func NewMemCache() *inMemcache {
	return &inMemcache{
		c:     make(map[string][]byte),
		mutex: sync.RWMutex{},
		Stat:  Stat{},
	}
}
