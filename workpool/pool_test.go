package workpool

import (
	"log"
	"testing"
)

func TestPool(t *testing.T) {
	task := NewTask(func() error {
		log.Println("hello")
		return nil
	})

	// p := NewPool(3, 0)
	// p := NewPool(3, 0, 100)
	// p := NewPool(3, 100)
	p := NewPool(3, 100, 100)

	go func() {
		i := 0
		for {
			//if i > 100000 {
			//	break
			//}

			p.AddTask(task)
			//log.Println("i = ", i)
			i++
		}
	}()

	p.Run()
}

func TestShutdown(t *testing.T) {
	task := NewTask(func() error {
		log.Println("hello")
		return nil
	})

	p := NewPool(10, 10000, 10000)

	go func() {
		i := 0
		for {
			if i > 1200000 {
				p.Shutdown()
				break
			}

			p.AddTask(task)
			log.Println("i = ", i)
			i++
		}
	}()

	p.Run()
}

/**
go test -v  -test.run TestPool
020/05/23 19:41:41 hello
2020/05/23 19:41:41 hello
2020/05/23 19:41:41 hello
2020/05/23 19:41:41 hello
2020/05/23 19:41:41 hello
2020/05/23 19:41:46 current worker id:  3 will exit...
2020/05/23 19:41:46 current worker id:  1 will exit...
2020/05/23 19:41:46 work pool shutdown success
2020/05/23 19:41:46 current worker id:  2 will exit...
--- PASS: TestPool (9.31s)
PASS
ok  	github.com/daheige/workpool	9.317s

go test -v  -test.run TestShutdown
2020/05/23 19:47:15 i =  1199999
2020/05/23 19:47:15 i =  1200000
2020/05/23 19:47:15 hello
2020/05/23 19:47:15 hello
2020/05/23 19:47:15 hello
2020/05/23 19:47:15 hello
2020/05/23 19:47:15 hello
2020/05/23 19:47:20 work pool will shutdown...
2020/05/23 19:47:20 recv signal:  terminated
2020/05/23 19:47:25 current worker id:  1 will exit...
2020/05/23 19:47:25 current worker id:  8 will exit...
2020/05/23 19:47:25 current worker id:  9 will exit...
2020/05/23 19:47:25 work pool shutdown success
2020/05/23 19:47:25 current worker id:  3 will exit...
--- PASS: TestShutdown (15.67s)
2020/05/23 19:47:25 current worker id:  5 will exit...
2020/05/23 19:47:25 current worker id:  2 will exit...
2020/05/23 19:47:25 current worker id:  10 will exit...
2020/05/23 19:47:25 current worker id:  6 will exit...
2020/05/23 19:47:25 current worker id:  7 will exit...
2020/05/23 19:47:25 current worker id:  4 will exit...
PASS
ok  	github.com/daheige/workpool	15.681s
PASS
ok  	github.com/daheige/workpool	15.354s
*/
