package jsonTime

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	type Person struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Birthday Time   `json:"birthday"`
	}

	// NullToEmptyStr = true

	now := Time(time.Now())
	t.Log(now)
	src := `{"id":5,"name":"xiaoming","birthday":"null"}`
	p := new(Person)
	err := json.Unmarshal([]byte(src), p)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(p)
	t.Log(time.Time(p.Birthday))
	js, _ := json.Marshal(p)
	t.Log(string(js))

	p2 := Person{
		Id:   1,
		Name: "fefe",
	}

	b, _ := json.Marshal(p2)

	log.Println("str: ", string(b))
}

/**
$ go test -v
=== RUN   TestTime
2019/10/12 23:07:52 str: {"id":1,"name":"fefe","birthday":null}
--- PASS: TestTime (0.00s)
    json_time_test.go:20: 2019-10-12 23:07:52
    json_time_test.go:28: &{5 xiaoming 0001-01-01 00:00:00}
    json_time_test.go:29: 0001-01-01 00:00:00 +0000 UTC
    json_time_test.go:31: {"id":5,"name":"xiaoming","birthday":null}
PASS
ok  	github.com/daheige/thinkgo/jsonTime	0.003s
*/
