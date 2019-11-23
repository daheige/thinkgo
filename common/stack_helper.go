package common

import (
	"log"
	"runtime/debug"
)

// CheckPanic check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		log.Println("panic error: ", err)
		log.Println(string(CatchStack()))
	}
}

// CatchStack 捕获指定stack信息,一般在处理panic/recover中处理
//返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	return debug.Stack()
}
