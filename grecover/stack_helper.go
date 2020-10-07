// Package grecover catch exec panic.
package grecover

import (
	"log"
	"os"
	"runtime/debug"
)

// Logger log interface
type Logger interface {
	Println(msg ...interface{})
}

// LogEntry log entry.
var (
	LogEntry   Logger = log.New(os.Stderr, "", log.LstdFlags)
	TracePanic        = false // trace panic stack
)

// CheckPanic check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		LogEntry.Println("panic error: ", err)
		if TracePanic {
			LogEntry.Println(string(CatchStack()))
		}
	}
}

// CatchStack 捕获指定stack信息,一般在处理panic/recover中处理
// 返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	return debug.Stack()
}
