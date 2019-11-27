package mytest

import (
	"encoding/json"
	"log"
	"testing"
)

/**json解析和反解析效率
$ time go test -v -test.run TestJson
2019/11/27 23:21:31 {"code":"200","count":49975004,"data":["golang","php","nodejs"],"message":"ok"}
2019/11/27 23:21:31 {"code":"200","count":49985002,"data":["golang","php","nodejs"],"message":"ok"}
2019/11/27 23:21:31 {"code":"200","count":49995001,"data":["golang","php","nodejs"],"message":"ok"}
--- PASS: TestJson (0.31s)
    json_test.go:19: start test json
PASS
ok  	github.com/daheige/thinkgo/mytest	0.317s

real	0m5.080s
user	0m1.373s
sys	0m0.422s
*/
func TestJson(t *testing.T) {
	t.Log("start test json")
	var data = map[string]interface{}{
		"code":    "200",
		"message": "ok",
		"data": []string{
			"golang", "php", "nodejs",
		},
		"count": 1,
	}

	for i := 0; i < 10000; i++ {
		data["count"] = data["count"].(int) + i
		json_str, _ := json.Marshal(data)
		log.Println(string(json_str))
	}
}
