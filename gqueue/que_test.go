package gqueue

import (
	"log"
	"testing"
)

func TestQue(t *testing.T) {

	taskCnt := 1000
	q := New(10, taskCnt)

	for i := 0; i < taskCnt; i++ {
		index := i

		q.Add(task(index))
	}

	log.Println(q)

	q.SetTaskCallback(taskCallback)
	q.SetFinishedCallback(taskFinished)
	q.Start()

}

func taskCallback(res interface{}) {
	log.Printf("current task result: %v", res)
}

func taskFinished() {
	log.Println("all task has finished")
}

func task(i int) func() interface{} {
	return func() interface{} {
		log.Println("exec task: ", i)
		return i
	}
}

/**
 * $ go test -v -test.run TestQue
 * 2018/10/28 15:06:00 all task has finished
--- PASS: TestQue (0.02s)
PASS
ok  	thinkgo/gqueue	0.020s
*/
