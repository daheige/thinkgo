package work

import (
	"fmt"
	"runtime/debug"
	"testing"
	"time"
)

type Score struct {
	Num int
}

//实现了job Do接口
func (s *Score) Do() {
	fmt.Println("num:", s.Num)
	time.Sleep(1 * 1 * time.Second)
}

func TestJobQueue(t *testing.T) {
	num := 100 * 100 * 2
	debug.SetMaxThreads(num) //设置最大线程数2w

	p := New(num)
	//模拟大量的job任务
	taskNums := 100 * 100 * 100
	go func() {
		for i := 1; i <= taskNums; i++ {
			sc := &Score{Num: i}
			p.Add(sc)
		}
	}()

	p.Wait()
}

/**
--- PASS: TestJobQueue (0.24s)
PASS
ok      github.com/daheige/thinkgo/work 0.734s
*/
