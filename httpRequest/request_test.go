package httpRequest

import (
	"log"
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	//请求句柄
	s := Service{
		BaseUri: "",
		Timeout: 2 * time.Second,
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
		return
	}

	//log.Println("data: ", string(res.Body))

	data := &ApiStdRes{}
	err := res.Json(data)
	log.Println(err)
	log.Println(data.Code, data.Message)
	log.Println(data.Data)

	res = s.Do("post", "http://localhost:1338/v1/data", &ReqOpt{
		Data: map[string]interface{}{
			"id": "1234",
		},
	})
	if res.Err != nil {
		log.Println("err: ", res.Err)
		return
	}

	log.Println(res.Err, string(res.Body))
}

/**
$ go test -v
2019/08/31 15:25:10 <nil> {"code":0,"data":["js","php","hello"],"message":"ok"}
--- PASS: TestRequest (0.25s)
PASS
*/
