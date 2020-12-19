package yamlconf

import (
	"log"
	"testing"
	"time"

	"github.com/daheige/thinkgo/gredigo"
	"github.com/daheige/thinkgo/mysql"
)

// Data test data.
type Data struct {
	gredigo.RedisConf
	Ip []string
}

func TestYaml(t *testing.T) {
	conf := NewConf()
	err := conf.LoadConf("test.yaml")
	log.Println(conf.GetData(), err)

	data := conf.GetData()

	var graceful time.Duration
	conf.Get("GracefulWait", &graceful)
	log.Println("graceful: ", graceful)

	log.Println("RedisCommon: ", data["RedisCommon"])

	// 读取数据到结构体中
	var redisConf = &Data{}
	conf.GetStruct("RedisCommon", redisConf)
	log.Println(redisConf)
	log.Println("Ip:", redisConf.Ip)
	log.Println(redisConf.Password == "")

	dbConf := &mysql.DbConf{}
	conf.GetStruct("DbDefault", dbConf)
	log.Println(dbConf)
}

/**
 * TestYaml
$ go test -v
=== RUN   TestYaml
2019/11/05 23:19:31 map[AppEnv:local AppName:hg-mux DbDefault:map[Database:test Ip:127.0.0.1 MaxIdleConns:10 MaxOpenConns:100 ParseTime:true Password:root Port:3306 ShowSql:true UsePool:true User:root] RedisCommon:map[Database:0 Host:127.0.0.1 IdleTimeout:120 Ip:[11.12.1.1 11.12.1.2 11.12.1.3] MaxActive:10 MaxIdle:3 Password:<nil> Port:6379]]
2019/11/05 23:19:31 RedisCommon:  map[Database:0 Host:127.0.0.1 IdleTimeout:120 Ip:[11.12.1.1 11.12.1.2 11.12.1.3] MaxActive:10 MaxIdle:3 Password:<nil> Port:6379]
2019/11/05 23:19:31 &{{127.0.0.1 6379  0 3 10 120} [11.12.1.1 11.12.1.2 11.12.1.3]}
2019/11/05 23:19:31 Ip: [11.12.1.1 11.12.1.2 11.12.1.3]
2019/11/05 23:19:31 true
2019/11/05 23:19:31 &{127.0.0.1 3306 root root test   10 100 true   <nil> true true}
--- PASS: TestYaml (0.00s)
PASS
ok  	github.com/daheige/tigago/yamlconf	0.005s
*/
