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
	// TraceFileLine(false) //关闭文件名和行数追踪

	SetLogDir("/web/wwwlogs/ilog")

	var nums int = 6 * 10000 //6w个独立协程处理,处理60w条日志写入

	var wg sync.WaitGroup
	wg.Add(nums) //一次性计数器设置，保证独立携程都成处理完毕

	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			Info("111222", map[string]interface{}{
				"id":   1234,
				"user": "heige",
			})

			Debug("this is debug: 111222", map[string]interface{}{
				"id":   12,
				"user": "daheige",
			})

			Error("error msg: 111222", nil)
			Notice("notice msg: 111222", nil)

			Warn("warning: 111222", map[string]interface{}{
				"name": "hello",
			})
			Critical("crit msg: 111222", nil)

			Alter("alter: 111222", nil)
			Emergency("emerg msg: 111222", nil)

			Alter("alter: 111222", nil)
			Emergency("emerg msg: 111222", nil)
		}()
	}

	wg.Wait()

	log.Println("write log success")

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	fmt.Println(now.Format(logTmMissMs)) //转换为Y-m-d H:i:s

	//获取文件名
	fmt.Println(filepath.Base("/mygo/src/thinkgo/common/Log.go"))
}

/**
$ time go test -v
$ time go test -v
=== RUN   TestLog
2019/12/01 11:51:37 write log success
2019-12-01 11:51:37
Log.go
--- PASS: TestLog (12.72s)
    log_test.go:13: 测试log库
PASS
ok  	github.com/daheige/thinkgo/glog	12.744s

real	0m14.047s
user	0m14.102s
sys	0m6.259s

qps: 47169 个/s 相比logger库基于zap封装的，qps要少1.2w左右
*/
