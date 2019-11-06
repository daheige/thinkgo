package rediscache

import (
	"time"

	"fmt"

	"github.com/gomodule/redigo/redis"
)

//redis连接信息
type RedisConf struct {
	Host        string
	Port        int
	Password    string
	Database    int
	MaxIdle     int //空闲pool个数
	MaxActive   int //最大激活数量
	IdleTimeout int //最大连接超时,单位s
}

var RedisPoolList = map[string]*redis.Pool{} //存放连接池信息

func GetRedisClient(name string) redis.Conn {
	if value, ok := RedisPoolList[name]; ok {
		return value.Get()
	}

	return nil
}

//添加新的redis连接池
func AddRedisPool(name string, conf *RedisConf) {
	RedisPoolList[name] = NewRedisPool(conf)
}

func (this *RedisConf) SetRedisPool(name string) {
	AddRedisPool(name, this)
}

//创建redis pool连接池
// If Wait is true and the pool is at the MaxActive limit, then Get() waits
// for a connection to be returned to the pool before returning.
func NewRedisPool(conf *RedisConf) *redis.Pool {
	return &redis.Pool{
		Wait:        true, //等待redis connection放入pool池子中
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		MaxActive:   conf.MaxActive,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port))
			if err != nil {
				return nil, err
			}

			if len(conf.Password) != 0 {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}

			//选择db
			if conf.Database >= 0 {
				if _, err := c.Do("SELECT", conf.Database); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
