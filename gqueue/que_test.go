package gqueue

import (
	"log"
	"testing"
)

func TestQue(t *testing.T) {
	q := New(10, 100)

	for i := 0; i < 1001; i++ {
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
