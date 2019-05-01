package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	SetLogFile("./logs/zap.log") //设置日志文件路径
	MaxSize(2)

	InitLogger()

	logSugar := LogSugar()
	logSugar.Debug(111)
	logSugar.Info(222)
	logSugar.Infof("hello,%s", "world")

	Info("111", zap.String("name", "abc"), zap.Int("age", 28))


	for i := 0; i < 100; i++ {
		Info("hello,world", zap.Int("a", 1), zap.String("b", "c"))
		Warn("haha")
	}

	Info("write success")
	Debug("hello")
	DPanic("111")

}
