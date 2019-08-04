package common

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"
)

// CheckPanic check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		loc, _ := time.LoadLocation("Local")
		fmt.Fprintf(os.Stderr, "\n%s %+v\n", FormatTime19(time.Now().In(loc)), err)
		fmt.Fprintf(os.Stderr, "full stack info: \n%s", CatchStack())
	}
}

// CatchStack 捕获指定stack信息,一般在处理panic/recover中处理
//返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	buf := &bytes.Buffer{}

	//完整的堆栈信息
	stack := Stack()
	buf.WriteString("full stack:\n")
	buf.Write(stack)

	//完整的函数调用信息
	buf.WriteString("full fn call info:\n")

	// skip为0时，打印当前调用文件及行数。
	// 为1时，打印上级调用的文件及行数，依次类推
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc).Name()
		buf.WriteString(fmt.Sprintf("error Stack file: %s:%d call func:%s\n", file, line, fn))
	}

	return buf.Bytes()
}

// Stack 获取完整的堆栈信息
// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false) //当第二个参数为true，一次获取所有的堆栈信息
		if n < len(buf) {
			return buf[:n]
		}

		buf = make([]byte, 2*len(buf))
	}
}
