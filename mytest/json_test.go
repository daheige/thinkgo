package mytest

import (
	"encoding/json"
	"log"
	"testing"
)

/**json解析和反解析效率
$ time go test -v -test.run TestJson
PASS
ok      github.com/daheige/thinkgo/mytest       0.718s

real    0m1.798s
user    0m1.492s
sys     0m0.208s
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
