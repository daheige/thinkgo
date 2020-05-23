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
2020/05/23 20:30:05 正在执行任务19997
[runner] 2020/05/23 20:30:05 current run task id:  19998
2020/05/23 20:30:05 正在执行任务19998
[runner] 2020/05/23 20:30:05 current run task id:  19999
2020/05/23 20:30:05 正在执行任务19999
[runner] 2020/05/23 20:30:05 task complete status:  <nil>
2020/05/23 20:30:05 error:  <nil>
2020/05/23 20:30:05 last_id:  19999
2020/05/23 20:30:05 all error:  map[]
--- PASS: TestRunner (1.19s)
PASS
ok      github.com/daheige/thinkgo/runner       1.201s
*/
