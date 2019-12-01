package logger

import (
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	SetLogDir("./logs/") //设置日志文件目录
	SetLogFile("mytest.log")
	MaxSize(20)
	TraceFileLine(false) //关闭文件名和行数追踪

	InitLogger(1)

	logSugar := LogSugar()
	logSugar.Debug(111)
	logSugar.Info(222)
	logSugar.Infof("hello,%s", "world")

	Info("111", map[string]interface{}{
		"abc": "daheige",
		"age": 28,
	})

	//测试60w日志输出到文件
	nums := 30 * 10000
	var wg sync.WaitGroup
	wg.Add(nums)
	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			Info("hello,world", map[string]interface{}{
				"a": 1,
				"b": "free",
			})

			Warn("haha", nil)
		}()
	}

	wg.Wait()

	Info("write success", nil)
	Error("type error", nil)
	Debug("hello", nil)
	DPanic("111", nil)
}

/**
$ time go test -v
=== RUN   TestLog
2019/12/01 11:47:07 msg:  hello
2019/12/01 11:47:07 log fields:  map[]
--- PASS: TestLog (10.75s)
PASS
ok  	github.com/daheige/thinkgo/logger	10.905s

real	0m31.969s
user	0m32.712s
sys	0m5.386s
qps: 55814 个/s

关闭记录日志文件名称和行号
$ time go test -v
=== RUN   TestLog
2019/12/01 11:49:06 msg:  hello
2019/12/01 11:49:06 log fields:  map[]
--- PASS: TestLog (10.14s)
PASS
ok  	github.com/daheige/thinkgo/logger	10.311s

real	0m11.933s
user	0m30.173s
sys	0m4.584s

qps: 59171 个/s
*/
