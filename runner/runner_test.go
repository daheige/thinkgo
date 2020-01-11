package runner

import (
	"log"
	"os"
	"testing"
	"time"
)

// TestRunner test runner
func TestRunner(t *testing.T) {

	// 日志句柄
	std := log.New(os.Stdout, "[runner] ", log.LstdFlags)

	// New参数可选，默认创建无超时的任务
	// p := New()

	// p := New(WithLogger(std))

	p := New(WithTimeout(3000*time.Millisecond), WithLogger(std))

	for i := 0; i < 20000; i++ {
		p.Add(createTask(i))
	}

	err := p.Start()
	log.Println("error: ", err)

	log.Println("last_id: ", p.GetLastTaskId())
	log.Println("all error: ", p.GetAllErrors())
}

// createTask 创建任务
func createTask(id int) func() error {
	return func() error {
		//panic(1)

		log.Printf("正在执行任务%d", id)
		// time.Sleep(time.Duration(id) * time.Millisecond)
		return nil
	}
}

/**
[runner] 2020/01/11 12:55:17 current run task id:  2050
2020/01/11 12:55:17 正在执行任务2050
[runner] 2020/01/11 12:55:17 received signal:  interrupt
[runner] 2020/01/11 12:55:17 task complete status:  received interrupt signal
2020/01/11 12:55:17 error:  received interrupt signal
2020/01/11 12:55:17 last_id:  2051
2020/01/11 12:55:17 all error:  map[]
--- PASS: TestRunner (1.16s)
PASS
ok      github.com/daheige/thinkgo/runner       2.680s
*/
