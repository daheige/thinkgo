package common

import (
	"time"

	"fmt"

	"github.com/gomodule/redigo/redis"
)

//redis连接信息
type RedisConf struct {
	Host        string
	Port        uint16
	Password    string
	Database    uint8
	MaxIdle     int   //空闲pool个数
	IdleTimeout int64 //最大连接超时,单位s
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
func NewRedisPool(conf *RedisConf) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
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
