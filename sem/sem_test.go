package sem

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestSemWithTimeout(t *testing.T) {
	tickets, timeout := 1, 3*time.Second
	s := New(tickets, timeout)

	if err := s.Acquire(); err != nil {
		panic(err)
	}

	// Do important work
	if err := s.Release(); err != nil {
		panic(err)
	}

	var cnt int
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int, wg *sync.WaitGroup) {
			defer func() {
				// 释放锁
				if err := s.Release(); err != nil {
					log.Println("release a sem err: ", err)
				}

				wg.Done()
			}()

			// 枷锁
			if err := s.Acquire(); err != nil {
				log.Println("get a sem err: ", err)
				return
			}

			log.Println("hello,world: ", i)
			cnt += i
		}(i, &wg)
	}

	wg.Wait()

	log.Println("cnt: ", cnt)

	log.Println("ok")

	// mutex lock without timeout
	// it will throw err
	t2, timeout2 := 0, 0
	s2 := New(t2, time.Duration(timeout2)*time.Second)

	if err := s2.Acquire(); err != nil {
		if err != ErrNoTickets {
			panic(err)
		}

		// No tickets left, can't work :(
		os.Exit(1)
	}
}

/**
$ go test -v
2019/08/04 11:13:27 hello,world:  716
2019/08/04 11:13:27 hello,world:  717
2019/08/04 11:13:27 hello,world:  704
2019/08/04 11:13:27 hello,world:  705
2019/08/04 11:13:27 hello,world:  708
2019/08/04 11:13:27 hello,world:  696
2019/08/04 11:13:27 cnt:  499500
2019/08/04 11:13:27 ok
exit status 1
FAIL	_/web/go/sem	0.016s
*/
