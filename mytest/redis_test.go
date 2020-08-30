package mytest

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/daheige/thinkgo/gredigo"

	"github.com/gomodule/redigo/redis"
)

func TestRedisPool(t *testing.T) {
	conf := &gredigo.RedisConf{
		Host:        "127.0.0.1",
		Port:        6379,
		MaxIdle:     100,
		MaxActive:   200,
		IdleTimeout: 240,
	}

	// 建立连接
	conf.SetRedisPool("default")
	var wg sync.WaitGroup

	for i := 0; i < 20000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := gredigo.GetRedisClient("default")
			defer client.Close()

			ok, err := client.Do("set", "myname", "daheige")
			fmt.Println(ok, err)

			value, _ := redis.String(client.Do("get", "myname"))
			fmt.Println("myname:", value)

			// 切换到database 1上面操作
			v, err := client.Do("Select", 1)
			fmt.Println(v, err)
			_, _ = client.Do("lpush", "myList", 123)
		}()
	}

	wg.Wait()
	log.Println("exec success...")

}

/*
$ go test -v -test.run TestRedisPool
OK <nil>
2019/11/06 22:07:50 exec success...
--- PASS: TestRedisPool (2.25s)
PASS
ok  	github.com/daheige/thinkgo/mytest	2.260s
*/
