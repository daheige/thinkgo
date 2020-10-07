package xerrors

import (
	"log"
	"testing"
)

func TestError(t *testing.T) {
	e := New("this is a error", 0, true)
	log.Println("error: ", e.Error())

	// 类型断言
	err := e.(*ErrorString)
	err.code = 123
	log.Println(err.Error())
	log.Println("error code: ", err.code)
	log.Println("full stack: ", string(err.Stack()))
	log.Printf("str: %+v", e) // 会调用Error()
}

/**
 * $ go test -v
=== RUN   TestError
2019/04/01 23:42:10 str: this is a error
--- PASS: TestError (0.00s)
PASS
ok  	github.com/daheige/thinkgo/xerrors	0.003
*/
