package mytest

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	fmt.Println("time: ", time.Now().Unix())                 //获取当前秒
	fmt.Println("ns time: ", time.Now().UnixNano())          //获取当前纳秒
	fmt.Println("ms time: ", time.Now().UnixNano()/1e6)      //将纳秒转换为毫秒
	fmt.Println("ns -> s time: ", time.Now().UnixNano()/1e9) //将纳秒转换为秒
	c := time.Unix(time.Now().UnixNano()/1e9, 0)             //将毫秒转换为 time 类型
	fmt.Println(c.String())                                  //输出当前英文时间戳格式

}

/**
$ go test -v -test.run TestTime
=== RUN   TestTime
time:  1534258904
ns time:  1534258904910270816
ms time:  1534258904910
ns -> s time:  1534258904
2018-08-14 23:01:44 +0800 CST
--- PASS: TestTime (0.00s)
PASS
ok      mytest  0.019s
*/
