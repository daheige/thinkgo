/**
自定义错误类型，一般用在api/微服务等业务逻辑中，处理错误
支持是否输出堆栈信息，可以把stack信息记录到日志文件中，方便定位问题
*/
package xerrors

import (
	"runtime"
)

type ErrorString struct {
	s     string
	Code  int
	frame []byte //错误堆栈信息
}

func New(text string, code int, isStack bool) error {
	return MakeError(text, code, isStack)
}

func MakeError(text string, code int, isStack bool) *ErrorString {
	err := &ErrorString{
		s:    text,
		Code: code,
	}

	if isStack {
		err.frame = stack()
	}

	return err
}

//实现了error interface{} Error方法
func (e *ErrorString) Error() string {
	return e.s
}

//打印完整的错误堆栈信息
func (e *ErrorString) Stack() []byte {
	return e.frame
}

//获取完整的堆栈信息
//捕获指定stack信息,一般在处理panic/recover中处理
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
