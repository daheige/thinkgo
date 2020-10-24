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
995: 07d5b6529ba84ac36304b034c3a5042a
996: 5c70a7511bf64b487f7bd697994af694
997: 1c6965d01aba48dd4aee7acb840650b1
998: d8ae9c92560f42144b80b6dfbfdc3693
999: 2fe7648ef2b748f55c2b67fa96dac2f6
--- PASS: TestUuid (0.02s)
PASS
*/
