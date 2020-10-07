package mytest

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/daheige/thinkgo/gutils"
)

func TestUuid(t *testing.T) {
	t.Log("测试uuid")
	fmt.Println(time.Now().UnixNano() / 1e6) // 将纳秒转换为毫秒
	fmt.Println(time.Now().UnixNano() / 1e9) // 将纳秒转换为毫秒
	fmt.Println(time.Now().Unix())           // 当前时间戳
	fmt.Println(1e3)                         // 1000

	var nums int = 1e3
	var strList []string

	for i := 0; i < nums; i++ {
		uuid := strings.Replace(gutils.NewUUID(), "-", "", -1)
		if checkExist(uuid, strList) {
			continue
		}

		strList = append(strList, uuid)
	}

	fmt.Println(len(strList))

	for k, v := range strList {
		fmt.Printf("%d: %s\n", k, v)
	}

}

func checkExist(str string, s []string) bool {
	for _, v := range s {
		if str == v {
			return true
		}
	}

	return false
}

/**
$ go test -v -test.run TestUuid
999: 2e750a8b5ae042f35cce573bb32eec96
--- PASS: TestUuid (0.01s)
    uuid_test.go:12: 测试uuid
PASS
ok      mytest  0.020s
*/
