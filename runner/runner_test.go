package runner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	log.Println("======开始执行任务=======")
	// timeout := 10 * time.Second //任务超时为3s
	// r := New(timeout)           //创建一个runner
	r := NewWithoutTime()
	for i := 0; i < 100; i++ {
		r.Add(createTask(i))
	}

	//开始执行任务
	if err := r.Start(); err != nil {
		switch err {
		case ErrInterrupt: //中断
			log.Println(err)
		case ErrorTimeout: //超时
			log.Println(err)
		}
	}

	log.Println("任务执行完毕")
	log.Println("错误的ids map: ", r.ErrorTaskIds())
	log.Println("最后执行的task id: ", r.GetLastTaskId())
}

func createTask(id int) func() error {
	return func() error {
		log.Printf("正在执行任务%d", id)
		time.Sleep(time.Duration(id) * time.Millisecond)
		return nil
	}
}

/**go test -v
2019/04/15 22:15:42 任务执行完毕
2019/04/15 22:15:42 错误的ids map:  map[]
2019/04/15 22:15:42 最后执行的task id:  99
--- PASS: TestRunner (4.98s)
PASS
ok  	github.com/daheige/thinkgo/runner	4.979s
*/
