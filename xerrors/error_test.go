package xerrors

import (
	"log"
	"testing"
)

func TestError(t *testing.T) {
	e := New("this is a error")
	// e.Code = 123

	log.Println(e.Error())
	var err error
	err = e
	log.Println(err.Error())
	// log.Println(e.Stack())
	log.Printf("str: %+v", e) //会调用Error()
}

/**
 * $ go test -v
=== RUN   TestError
2019/04/01 23:42:10 this is a error
2019/04/01 23:42:10 this is a error
2019/04/01 23:42:10 str: this is a error
--- PASS: TestError (0.00s)
PASS
ok  	github.com/daheige/thinkgo/xerrors	0.003
*/
