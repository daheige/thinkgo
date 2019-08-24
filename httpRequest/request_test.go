package httpRequest

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	//请求句柄
	s := Service{
		BaseUri: "",
		Timeout: 1,
	}

	//请求参数设置
	opt := &ReqOpt{
		Params: map[string]interface{}{
			"objid":   12784,
			"objtype": 1,
			"p":       0,
		},
	}

	res := s.Do("get", "https://studygolang.com/object/comments", opt)
	if res.Err != nil {
		log.Println("err: ", res.Err)
		t.Error(res.Err)
		return
	}

	//log.Println("data: ", string(res.Body))

	data := &ApiStdRes{}
	err := res.Json(data)
	log.Println(err)
	log.Println(data.Code, data.Message)
	log.Println(data.Data)

}

/**
$ go test -v
--- PASS: TestRequest (0.26s)
PASS
ok      github.com/daheige/httpRequest     0.265s
*/
