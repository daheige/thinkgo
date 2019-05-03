package logger

import (
	"sync"
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	SetLogFile("./logs/zap.log") //设置日志文件路径
	MaxSize(20)

	InitLogger()

	logSugar := LogSugar()
	logSugar.Debug(111)
	logSugar.Info(222)
	logSugar.Infof("hello,%s", "world")

	Info("111", zap.String("name", "abc"), zap.Int("age", 28))

	nums := 30 * 10000
	var wg sync.WaitGroup
	wg.Add(nums)
	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			Info("hello,world", zap.Int("a", 1), zap.String("b", "c"))
			Warn("haha")
		}()
	}

	wg.Wait()

	Info("write success")
	Debug("hello")
	DPanic("111")

}

/**
$ go test -v
=== RUN   TestLog
--- PASS: TestLog (10.38s)
PASS
ok  	logger	10.458s
 */
