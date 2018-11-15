package common

import (
	"log"
	"testing"
)

func TestYaml(t *testing.T) {
	conf := NewConf()
	conf.LoadConf("test.yaml")
	log.Println(conf.data)

	// var v interface{}
	// conf.GetStruct("RedisCommon", v)
	var v = &RedisConf{}
	conf.GetStruct("RedisCommon", v)
	log.Println(v)
	log.Println(conf.data["RedisCommon"])

}

/**
 * TestYaml
=== RUN   TestYaml
2018/11/15 20:39:05 map[AppEnv:local AppName:hg-mux RedisCommon:map[IdleTimeout:120 Host:127.0.0.1 Port:6379 Password:<nil> Database:0 MaxIdle:3 MaxActive:10]]
2018/11/15 20:39:05 map[interface {}]interface {}
2018/11/15 20:39:05 k =  MaxIdle
2018/11/15 20:39:05 k type = string
2018/11/15 20:39:05 v =  3
2018/11/15 20:39:05 k =  MaxActive
2018/11/15 20:39:05 k type = string
2018/11/15 20:39:05 v =  10
2018/11/15 20:39:05 k =  IdleTimeout
2018/11/15 20:39:05 k type = string
2018/11/15 20:39:05 v =  120
2018/11/15 20:39:05 k =  Host
2018/11/15 20:39:05 k type = string
2018/11/15 20:39:05 v =  127.0.0.1
2018/11/15 20:39:05 k =  Port
2018/11/15 20:39:05 k type = string
2018/11/15 20:39:05 v =  6379
2018/11/15 20:39:05 k =  Database
2018/11/15 20:39:05 k type = string
2018/11/15 20:39:05 v =  0
2018/11/15 20:39:05 &{127.0.0.1 6379  0 3 10 120}
2018/11/15 20:39:05 map[Database:0 MaxIdle:3 MaxActive:10 IdleTimeout:120 Host:127.0.0.1 Port:6379 Password:<nil>]
--- PASS: TestYaml (0.02s)
PASS
ok  	thinkgo/common	0.024s
*/
