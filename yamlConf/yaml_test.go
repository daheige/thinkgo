package yamlConf

import (
	"log"
	"testing"
	"github.com/daheige/thinkgo/redisCache"
)

func TestYaml(t *testing.T) {
	conf := NewConf()
	conf.LoadConf("test.yaml")
	log.Println(conf.data)

	// var v interface{}
	// conf.GetStruct("RedisCommon", v)
	var v = &redisCache.RedisConf{}
	conf.GetStruct("RedisCommon", v)
	log.Println(v)
	log.Println(conf.data["RedisCommon"])

}

/**
 * TestYaml
$ go test -v
=== RUN   TestYaml
2019/03/17 20:29:13 map[AppEnv:local AppName:hg-mux RedisCommon:map[Database:0 Host:127.0.0.1 IdleTimeout:120 MaxActive:10 MaxIdle:3 Password:<nil> Port:6379]]
2019/03/17 20:29:13 &{127.0.0.1 6379  0 3 10 120}
2019/03/17 20:29:13 map[Database:0 Host:127.0.0.1 IdleTimeout:120 MaxActive:10 MaxIdle:3 Password:<nil> Port:6379]
--- PASS: TestYaml (0.00s)
PASS
ok      github.com/daheige/thinkgo/yamlConf     0.005s
*/
