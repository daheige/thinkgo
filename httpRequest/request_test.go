package httpRequest

import (
	"log"
	"testing"
)

func TestRequest(t *testing.T) {
	r := ApiRequest{
		BaseUri: "",
		Url:     "https://studygolang.com/object/comments?objid=12784",
		Params: map[string]interface{}{
			// "objid": 12784,
			"objtype": 1,
			"p":       0,
		},
		Method: "get",
	}

	res := r.Do()
	if res.Err != nil {
		log.Println("err: ", res.Err)
		t.Error(res.Err)
		return
	}

	log.Println("data: ", res.Body)

}

/**
$ go test -v
=== RUN   TestRequest
2019/03/30 17:47:56 data:  xxx
--- PASS: TestRequest (0.26s)
PASS
ok  	httpRequest	0.265s
 */
