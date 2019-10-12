package gResty

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
		RetryCount: 2,
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
		RetryCount: 2, //重试次数
	})
	if res.Err != nil {
		log.Println("err: ", res.Err)
		return
	}

	log.Println(res.Err, string(res.Body))
}

/**
$ go test -v
2019/10/12 21:38:56 <nil> {"code":0,"data":["js","php",""],"message":"ok"}
--- PASS: TestRequest (0.41s)
PASS
*/
