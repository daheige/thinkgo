/**
自定义错误类型，一般用在api/微服务等业务逻辑中，处理错误
支持是否输出堆栈信息，可以把stack信息记录到日志文件中，方便定位问题
*/
package xerrors

import "runtime/debug"

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
		err.frame = debug.Stack()
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
