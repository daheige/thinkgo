package mytest

import (
	"fmt"
	"testing"
	"thinkgo/common"
)

type NginxConf struct {
	Host    string
	Port    int
	MaxFile int
}

//验证common yaml读取配置文件
func TestYamlRead(t *testing.T) {
	t.Log("测试yaml读取配置文件")

	conf := common.New()
	path := "test.yaml"
	conf.LoadConf(path)
	fmt.Println("读取的内容:", conf.Data)

	fmt.Println(conf.GetString("App_key", ""))
	fmt.Println(conf.GetString("App_key", ""))

	//读取NginxConf
	nginxConf := &NginxConf{}
	fmt.Println(conf.GetStruct("NginxConf", nginxConf))
	fmt.Println(nginxConf)

	/*//不存在的key指定默认值1111
	  fmt.Println(conf.GetString("App_key2", "1111"))

	  //int ---> int64
	  a := 12342343
	  fmt.Println(int64(a))

	  fmt.Println(strconv.ParseInt("1234", 10, 64))
	  fmt.Println(conf.GetInt64("Count", 1234))
	  fmt.Println(conf.GetFloat64("Count2", 1234.12))
	  fmt.Println(conf.GetString("Redis.Default.Host", "")) //支持.的方式读取内容
	  fmt.Println(conf.GetInt("NginxConf.MaxFile", 12342))  //支持.的方式读取内容*/

}

/**
$ go test -v -test.run TestYamlRead
=== RUN   TestYamlRead
读取的内容: map[Database:map[Default:map[Port:3306 User:root Pwd:1234 Database:test Host:127.0.0.1] Default2:map[Host:127.0.0.1 Port:3306 User:root Pwd:1234 Database:test]] NginxConf:map[MaxFile:65355 Host:127.0.0.1 Port:3309] Redis:map[Default:map[Host:127.0.0.1 Port:6379] Order:map[Host:127.0.0.1 Port:6379]] App_key:1234fessfe App_env:testing Count:1234]
1234fessfe
1234fessfe
&{127.0.0.1 3309 65355}
&{127.0.0.1 3309 65355}
1111
12342343
1234 <nil>
1234
1234.12
127.0.0.1
65355
--- PASS: TestYamlRead (0.00s)
    yaml_test.go:18: 测试yaml读取配置文件
PASS
ok      mytest  0.002s

//建议定义结构体的方式读取
$ go test -v
=== RUN   TestYamlRead
读取的内容: map[Count:1234 Database:map[Default:map[Pwd:1234 Database:test Host:127.0.0.1 Port:3306 User:root] Default2:map[Host:127.0.0.1 Port:3306 User:root Pwd:1234 Database:test]] NginxConf:map[Host:127.0.0.1 Port:3309 MaxFile:65355] Redis:map[Default:map[Host:127.0.0.1 Port:6379] Order:map[Host:127.0.0.1 Port:6379]] App_key:1234fessfe App_env:testing]
1234fessfe
1234fessfe
&{127.0.0.1 3309 65355}
&{127.0.0.1 3309 65355}
--- PASS: TestYamlRead (0.00s)
    yaml_test.go:17: 测试yaml读取配置文件
PASS
ok      mytest  0.004s
*/
