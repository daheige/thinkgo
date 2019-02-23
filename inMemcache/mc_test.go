package inMemcache

import (
	"log"
	"testing"
)

func TestMemCache(t *testing.T) {
	log.Println("start test...")
	var c Cacher
	inCache := NewMemCache()

	c = inCache

	//接口变量包含了两个部分，接口对于的具体类型和该类型的值
	log.Printf("c type is %T,value is %v", c, c)
	//c type is *cache.inMemcache,value is &{map[] {{0 0} 0 0 0 0} {0 0 0}}

	c.Set("daheige", []byte("12fff"))

	c.Set("k2", []byte("12fff"))

	c.Set("k2", []byte("12345"))

	c.Set("k3", []byte("golang"))

	v, err := c.Get("daheige")
	if err != nil {
		log.Println(err.Error())
		return
	}

	s := c.GetStat()
	log.Println(s.Count)
	log.Println(s.KeySize)
	log.Println(s.ValueSize)

	log.Println(ToString(c.Get("k3")))
	log.Println(string(v))
}

/*
$ go test -v
=== RUN   TestMemCache
2019/02/23 17:55:14 start test...
2019/02/23 17:55:14 c type is *inMemcache.inMemcache,value is &{map[] {{0 0} 0 0 0 0} {0 0 0}}
2019/02/23 17:55:14 3
2019/02/23 17:55:14 11
2019/02/23 17:55:14 16
2019/02/23 17:55:14 golang <nil>
2019/02/23 17:55:14 12fff
--- PASS: TestMemCache (0.00s)
PASS
ok  	github.com/daheige/thinkgo/inMemcache	0.002s
*/
