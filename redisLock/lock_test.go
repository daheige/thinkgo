package redisLock

import (
	"log"
	"sync"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func lock() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println("redis connection error: ", err)
		return
	}

	defer conn.Close()

	l := New(conn, "heige", "hello,world", 100)

	if ok, err := l.TryLock(); ok {
		log.Println("lock success")
		for i := 0; i < 10; i++ {
			log.Println("hello,i: ", i)
		}

		l.Unlock()
	} else {
		log.Println("lock fail")
		log.Println("err: ", err)
	}
}

//测试枷锁操作
func TestRedisLock(t *testing.T) {
	lock()
}

//并发操作的尝试枷锁
func TestLock(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func(wg *sync.WaitGroup) {
			defer wg.Done()

			lock()

		}(&wg)
	}

	wg.Wait()
	log.Println("ok")
}

/**
2019/08/10 23:36:13 lock fail
2019/08/10 23:36:13 err:  <nil>
2019/08/10 23:36:13 lock fail
2019/08/10 23:36:13 err:  <nil>
2019/08/10 23:36:13 lock success
2019/08/10 23:36:13 hello,i:  0
2019/08/10 23:36:13 hello,i:  1
2019/08/10 23:36:13 hello,i:  2
2019/08/10 23:36:13 hello,i:  3
2019/08/10 23:36:13 hello,i:  4
2019/08/10 23:36:13 hello,i:  5
2019/08/10 23:36:13 hello,i:  6
2019/08/10 23:36:13 hello,i:  7
2019/08/10 23:36:13 lock fail
2019/08/10 23:36:13 err:  <nil>
2019/08/10 23:36:13 hello,i:  8
2019/08/10 23:36:13 hello,i:  9
2019/08/10 23:36:13 ok
--- PASS: TestLock (0.02s)
PASS
ok      redisLock       0.054s

*/
