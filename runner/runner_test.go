package runner

import (
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T) {
	log.Println("======开始执行任务=======")
	timeout := 10 * time.Second //任务超时为3s
	r := New(timeout)           //创建一个runner
	for i := 0; i < 1000; i++ {
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
	log.Println("已完成的任务ids: ", r.GetDoneTaskIds())
	log.Println("最后执行的task id: ", r.GetLastTaskId())
}

func createTask(id int) func() {
	return func() {
		log.Printf("正在执行任务%d", id)
		time.Sleep(time.Duration(id) * time.Millisecond)
	}
}

/**go test -v
2018/09/16 21:38:45 received timeout
2018/09/16 21:38:45 任务执行完毕
2018/09/16 21:38:45 已完成的任务ids:  [0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 41 42 43 44 45 46 47 48 49 50 51 52 53 54 55 56 57 58 59 60 61 62 63 64 65 66 67 68 69 70 71 72 73 74 75 76]
2018/09/16 21:38:45 最后执行的task id:  76
--- PASS: TestRunner (3.00s)
PASS
ok  	thinkgo/runner	3.003s
*/
