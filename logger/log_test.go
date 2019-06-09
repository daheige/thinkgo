package logger

import (
	"sync"
	"testing"
)

func TestLog(t *testing.T) {
	SetLogDir("./logs/") //设置日志文件目录
	SetLogFile("mytest.log")
	MaxSize(20)

	InitLogger()

	logSugar := LogSugar()
	logSugar.Debug(111)
	logSugar.Info(222)
	logSugar.Infof("hello,%s", "world")

	Info("111", map[string]interface{}{
		"abc": "daheige",
		"age": 28,
	})

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
$ go test -v
=== RUN   TestLog
2019/06/09 17:13:50 msg:  hello
2019/06/09 17:13:50 log fields:  map[]
--- PASS: TestLog (13.12s)
PASS
ok      github.com/daheige/thinkgo/logger       13.269s
*/
