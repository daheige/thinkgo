package goredis

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/daheige/thinkgo/common"
	"github.com/go-redis/redis"
)

func TestRedis(t *testing.T) {
	conf := RedisClientConf{
		Address:     "127.0.0.1:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		PoolSize:    10,
		PoolTimeout: 10 * time.Second,
	}

	// client := conf.GetClient()
	// defer client.Close()

	conf.SetClientName("default")

	client, err := GetRedisClient("default")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("username", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	//redis cluster test
	clusterConf := RedisClusterConf{
		AddressNodes: []string{
			"127.0.0.1:6391",
			"127.0.0.1:6392",
			"127.0.0.1:6393",
			"127.0.0.1:6394",
			"127.0.0.1:6395",
			"127.0.0.1:6396",
		},
		PoolSize:     10, // PoolSize applies per cluster node and not for the whole cluster.
		MaxRetries:   2,  //重试次数
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second, //底层默认3s
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
		MinIdleConns: 10,
		IdleTimeout:  100 * time.Second,
	}

	cluster := clusterConf.GetCluster()

	defer cluster.Close()

	str, err := cluster.Set("username", "daheige", 1000*time.Second).Result()
	log.Println(str, err)

	str, err = cluster.Set("myname", "daheige2", 1000*time.Second).Result()
	log.Println(str, err)

	log.Println(cluster.Get("username").Result()) //2019/09/26 21:42:18 daheige <nil>
}

/**
=== RUN   TestRedis
PONG <nil>
2019/11/06 22:11:08 OK <nil>
2019/11/06 22:11:08 OK <nil>
2019/11/06 22:11:08 daheige <nil>
--- PASS: TestRedis (0.03s)
PASS
ok  	github.com/daheige/thinkgo/goredis	0.028s
*/

// go test -v -test.run=TestRedis2
func TestRedis2(t *testing.T) {
	conf := RedisClientConf{
		Address:     "127.0.0.1:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		PoolSize:    10,
		PoolTimeout: 10 * time.Second,
	}

	conf.SetClientName("default")

	client, err := GetRedisClient("default")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)

		go setData(client, &wg)
	}

	wg.Wait()

	log.Println("ok")
}

func TestRedis3(t *testing.T) {
	conf := RedisClientConf{
		Address:     "127.0.0.1:6379",
		Password:    "", // no password set
		DB:          0,  // use default DB
		PoolSize:    10,
		PoolTimeout: 10 * time.Second,
	}

	conf.SetClientName("default")

	client, err := GetRedisClient("default")
	if err != nil {
		panic(err)
	}

	defer client.Close()

	// 通道乐观锁
	lock := common.NewChanLock()

	var wg sync.WaitGroup
	for i := 0; i < 5000; i++ {
		wg.Add(1)

		go setData2(client, &wg, lock)
	}

	wg.Wait()

	log.Println("ok")
}

func setData(client *redis.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	key := "mytest"

	str, _ := client.Do("get", key).String()
	// log.Println("err: ", err)

	if str != "" {
		log.Println("str: ", str)
		return
	}

	// 第一次运行的时候，会让多个goroutine走到设置缓存下面
	log.Println("start set redis data")

	client.Do("setEx", key, 10, "1111")
}

func setData2(client *redis.Client, wg *sync.WaitGroup, lock *common.ChanLock) {
	defer wg.Done()

	key := "mytest"

	str, _ := client.Do("get", key).String()
	// log.Println("err: ", err)

	if str != "" {
		log.Println("str: ", str)
		return
	}

	//采用乐观锁实现，但这种方式只在单机上才可以，如果是多个机器，就需要用分布式锁
	if lock.TryLock() {
		log.Println("start set redis data")

		client.Do("setEx", key, 10, "1111")
	}
}

/**
go test -v -test.run=TestRedis2
第一次运行的时候，会让多个goroutine走到设置缓存下面
=== RUN   TestRedis2
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 str:  1111
2019/11/20 21:00:36 str:  1111
2019/11/20 21:00:36 str:  1111
2019/11/20 21:00:36 str:  1111
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 start set redis data
2019/11/20 21:00:36 str:  1111
2019/11/20 21:00:36 str:  1111

第二次，执行,redis设置的次数减少
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 start set redis data
2019/11/20 21:01:58 str:  1111
2019/11/20 21:01:58 str:  1111
2019/11/20 21:01:58 str:  1111
2019/11/20 21:01:58 str:  1111

综合上面的测试，需要对redis设置的时候，只需要一次执行就可以
可以加上一个分布式锁或乐观锁实现
*/
