package mytest

import (
	"fmt"
	"testing"

	"thinkgo/common"

	"github.com/gomodule/redigo/redis"
)

func TestRedisPool(t *testing.T) {
	conf := &common.RedisConf{
		Host: "127.0.0.1",
		Port: 6379,
	}

	//建立连接
	conf.SetRedisPool("default")

	client := common.GetRedisClient("default")
	defer client.Close()

	ok, err := client.Do("set", "myname", "daheige")
	fmt.Println(ok, err)

	value, err := redis.String(client.Do("get", "myname"))
	fmt.Println("myname:", value)

	//切换到database 1上面操作
	v, err := client.Do("Select", 1)
	fmt.Println(v, err)
	client.Do("lpush", "myList", 123)

}
