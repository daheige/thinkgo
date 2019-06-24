package yamlConf

import (
	"log"
	"testing"

	"github.com/daheige/thinkgo/redisCache"
)

type Data struct {
	redisCache.RedisConf
	Ip []string
}

func TestYaml(t *testing.T) {
	conf := NewConf()
	conf.LoadConf("test.yaml")
	log.Println(conf.data)

	log.Println("RedisCommon: ", conf.data["RedisCommon"])

	//读取数据到结构体中
	var redisConf = &Data{}
	conf.GetStruct("RedisCommon", redisConf)
	log.Println(redisConf)
	log.Println("Ip:", redisConf.Ip)
	log.Println(redisConf.Password == "")
}

/**
 * TestYaml
$ go test -v
=== RUN   TestYaml
2019/06/24 23:10:03 map[AppEnv:local AppName:hg-mux RedisCommon:map[Database:0 Host:127.0.0.1 IdleTimeout:120 Ip:[11.12.1.1 11.12.1.2 11.12.1.3] MaxActive:10 MaxIdle:3 Password:<nil> Port:6379]]
2019/06/24 23:10:03 RedisCommon:  map[Database:0 Host:127.0.0.1 IdleTimeout:120 Ip:[11.12.1.1 11.12.1.2 11.12.1.3] MaxActive:10 MaxIdle:3 Password:<nil> Port:6379]
2019/06/24 23:10:03 s &{{ 0  0 0 0 0} []}
2019/06/24 23:10:03 &{{127.0.0.1 6379  0 3 10 120} [11.12.1.1 11.12.1.2 11.12.1.3]}
2019/06/24 23:10:03 Ip: [11.12.1.1 11.12.1.2 11.12.1.3]
2019/06/24 23:10:03 true
--- PASS: TestYaml (0.00s)
PASS
ok      github.com/daheige/thinkgo/yamlConf     0.005s
*/
