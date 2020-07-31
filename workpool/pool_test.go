package workpool

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	task := NewTask(func() error {
		log.Println("hello")
		return nil
	})

	p := NewPool(
		WithExecInterval(100*time.Millisecond),
		WithEntryCap(10), WithJobCap(10), WithWorkerCap(10),
		WithLogger(log.New(os.Stderr, "", log.LstdFlags)),
	)

	// p := NewPool(WithEntryCap(3), WithWorkerCap(10))

	go func() {
		i := 0
		for {
			// if i > 100000 {
			//	break
			// }

			p.AddTask(task)
			// log.Println("i = ", i)
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

	p := NewPool(
		WithExecInterval(100*time.Millisecond),
		WithEntryCap(10), WithJobCap(1000),
		WithWorkerCap(1000), WithLogger(log.New(os.Stderr, "", log.LstdFlags)),
		WithEntryCloseWait(2*time.Second),
		WithShutdownWait(3*time.Second),
	)

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
2020/07/04 11:33:01 i =  937710
2020/07/04 11:33:01 current worker id:  68 will exit...
2020/07/04 11:33:01 current worker id:  40 will exit...
2020/07/04 11:33:01 current worker id:  447 will exit...
2020/07/04 11:33:01 current worker id:  734 will exit...
2020/07/04 11:33:01 current worker id:  579 will exit...
2020/07/04 11:33:01 current worker id:  737 will exit...
2020/07/04 11:33:01 current worker id:  20 will exit...
2020/07/04 11:33:01 i =  937711
2020/07/04 11:33:01 i =  937712
2020/07/04 11:33:01 i =  937713
2020/07/04 11:33:01 current worker id:  995 will exit...
2020/07/04 11:33:01 work pool shutdown success
2020/07/04 11:33:01 i =  937719
2020/07/04 11:33:01 i =  937720
2020/07/04 11:33:01 i =  937721
--- PASS: TestShutdown (4.58s)
2020/07/04 11:33:01 i =  937731
PASS
ok  	github.com/daheige/workpool	4.586s

go test -v  -test.run TestShutdown
2020/07/04 11:16:07 hello
2020/07/04 11:16:07 current worker id:  457
2020/07/04 11:16:07 hello
2020/07/04 11:16:07 current worker id:  245
2020/07/04 11:16:09 work pool will shutdown...
2020/07/04 11:16:09 recv signal:  terminated
2020/07/04 11:16:11 current worker id:  740 will exit...
2020/07/04 11:16:11 current worker id:  245 will exit...
2020/07/04 11:16:11 current worker id:  1000 will exit...
2020/07/04 11:16:11 current worker id:  6 will exit...
2020/07/04 11:16:11 current worker id:  723 will exit...
2020/07/04 11:16:11 current worker id:  334 will exit...
2020/07/04 11:16:11 current worker id:  754 will exit...
2020/07/04 11:16:11 work pool shutdown success
2020/07/04 11:16:11 current worker id:  324 will exit...
2020/07/04 11:16:11 current worker id:  423 will exit...
--- PASS: TestShutdown (127.41s)
2020/07/04 11:16:11 current worker id:  830 will exit...
2020/07/04 11:1PASS
6:11 current worker id:  693 will exit...
2020/07/04 11:16:11 current worker id:  552 will exit...
2020/07/04 11:16:11 current worker id:  394 will exit...
2020/07/04 11:16:11 current worker id:  14 will exit...
*/
