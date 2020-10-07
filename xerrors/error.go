/**
自定义错误类型，一般用在api/微服务等业务逻辑中，处理错误
支持是否输出堆栈信息，可以把stack信息记录到日志文件中，方便定位问题
*/
package xerrors

import "runtime/debug"

type ErrorString struct {
	msg   string
	code  int
	frame []byte // 错误堆栈信息
}

// New 创建一个error
func New(text string, code int, isStack ...bool) error {
	var b bool
	if len(isStack) > 0 && isStack[0] {
		b = true
	}

	return MakeError(text, code, b)
}

// MakeError 创建一个error
func MakeError(text string, code int, isStack bool) *ErrorString {
	err := &ErrorString{
		msg:  text,
		code: code,
	}

	if isStack {
		err.frame = debug.Stack()
	}

	return err
}

// Error 实现了error interface{} Error方法
func (e *ErrorString) Error() string {
	return e.msg
}

// Code 返回code
func (e *ErrorString) Code() int {
	return e.code
}

// Stack 打印完整的错误堆栈信息
func (e *ErrorString) Stack() []byte {
	return e.frame
}
