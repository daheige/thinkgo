package xerrors

import (
	"bytes"
	"fmt"
	"runtime"
)

type ErrorString struct {
	s    string
	code int
}

func New(text string) *ErrorString {
	return &ErrorString{
		s: text,
	}
}

func (e *ErrorString) SetCode(code int) {
	e.code = code
}

func (e *ErrorString) ErrCode() int {
	return e.code
}

//实现了error interface{} Error方法
func (e *ErrorString) Error() string {
	return e.s
}

//打印完整的错误堆栈信息
func (e *ErrorString) Stack() string {
	return string(fullStack())
}

//捕获指定stack信息,一般在处理panic/recover中处理
//返回完整的堆栈信息和函数调用信息
func fullStack() []byte {
	buf := &bytes.Buffer{}

	//完整的堆栈信息
	stack := stack()
	buf.WriteString("full stack:\n")
	buf.WriteString(string(stack))

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

//获取完整的堆栈信息
// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false) //当第二个参数为true，一次获取所有的堆栈信息
		if n < len(buf) {
			return buf[:n]
		}

		buf = make([]byte, 2*len(buf))
	}
}
