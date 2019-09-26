package goRedis

import (
	"fmt"
	"log"
	"testing"
	"time"
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
2019/09/26 21:48:30 OK <nil>
2019/09/26 21:48:30 OK <nil>
2019/09/26 21:48:30 daheige <nil>
--- PASS: TestRedis (0.01s)
PASS
*/
