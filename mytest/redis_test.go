package mytest

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/daheige/thinkgo/redisCache"

	"github.com/gomodule/redigo/redis"
)

func TestRedisPool(t *testing.T) {
	conf := &redisCache.RedisConf{
		Host:        "127.0.0.1",
		Port:        6379,
		MaxIdle:     100,
		MaxActive:   200,
		IdleTimeout: 240,
	}

	//建立连接
	conf.SetRedisPool("default")
	var wg sync.WaitGroup

	for i := 0; i < 20000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := redisCache.GetRedisClient("default")
			defer client.Close()

			ok, err := client.Do("set", "myname", "daheige")
			fmt.Println(ok, err)

			value, err := redis.String(client.Do("get", "myname"))
			fmt.Println("myname:", value)

			//切换到database 1上面操作
			v, err := client.Do("Select", 1)
			fmt.Println(v, err)
			client.Do("lpush", "myList", 123)
		}()
	}

	wg.Wait()
	log.Println("exec success...")

}

/*
$ go test -v -test.run TestRedisPool
2019/03/20 22:56:04 exec success...
--- PASS: TestRedisPool (0.95s)
PASS
ok  	github.com/daheige/thinkgo/mytest	1.902s
*/
