package glog

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Log("测试log库")

	LogSplit(true)
	LogSize(30) //单个日志文件大小

	SetLogDir("/web/wwwlogs/ilog")
	var wg sync.WaitGroup
	var nums int = 10 * 10000 //10w个独立协程处理,处理100w条日志写入
	wg.Add(nums)              //一次性计数器设置，保证独立携程都成处理完毕

	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			InfoLog("111222", map[string]interface{}{
				"id":   1234,
				"user": "heige",
			})

			DebugLog("this is debug: 111222", map[string]interface{}{
				"id":   12,
				"user": "daheige",
			})

			ErrorLog("error msg: 111222", nil)
			NoticeLog("notice msg: 111222", nil)

			WarnLog("warning: 111222", map[string]interface{}{
				"name": "hello",
			})
			CritLog("crit msg: 111222", nil)

			AlterLog("alter: 111222", nil)
			EmergLog("emerg msg: 111222", nil)

			AlterLog("alter: 111222", nil)
			EmergLog("emerg msg: 111222", nil)
		}()
	}

	wg.Wait()

	log.Println("write log success")

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	fmt.Println(now.Format(logtmFmtTime)) //转换为Y-m-d H:i:s

	//获取文件名
	fmt.Println(filepath.Base("/mygo/src/thinkgo/common/Log.go"))
}

/**
$ go test -v -test.run=TestLog
=== RUN   TestLog
2019/06/09 13:28:21 write log success
2019-06-09
Log.go
--- PASS: TestLog (28.73s)
    log_test.go:13: 测试log库
PASS
ok      github.com/daheige/thinkgo/common       28.779s
*/
