package mytest

import (
	"encoding/json"
	"testing"

	myjson "github.com/daheige/thinkgo/jsoniter"
)

//测试github.com/json-iterator/go
//json解析和反解析效率
/**
$ time go test -v -test.run TestJson
PASS
ok      mytest  0.220s

real    0m1.015s
user    0m1.179s
sys 0m0.154s
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
		t.Log(string(json_str))
	}
}

/**
time go test -v test.run TestMyJson
PASS
ok      mytest  0.173s

real    0m0.898s
user    0m1.086s
sys 0m0.160s
*/
func TestMyJson(t *testing.T) {
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
		json_str, _ := myjson.Marshal(data)
		t.Log(string(json_str))
	}
}
