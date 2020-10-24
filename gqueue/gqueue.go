// package gqueue 通过指定goroutine个数,实现task queue执行器
// 提交任务到tash chan中，然后不断从chan中取出task执行
// 结合官方的sync.WaitGroup计数信号等待执行完毕
// go goroutine非抢占式的,通过runtime.Gosched()让出cpu给其他goroutine
package gqueue

import (
	"runtime"
	"sync"
)

type Queue struct {
	gNum             int                     // 并发执行任务所需要的goroutine个数
	taskTotal        int                     // 执行任务的总数
	tasks            chan func() interface{} // 任务放置在缓冲通道中
	taskCallback     func(res interface{})   // 每个任务执行后的回调函数
	finishedCallback func()                  // 所有任务执行完毕后的回调
	wg               sync.WaitGroup          // 保证goroutine同步执行的信号计数器
}

// New 创建一个任务队列实例
func New(number, total int) *Queue {
	if total < 1 {
		panic("task total number must gt 1")
	}

	if number > total {
		number = total
	}

	return &Queue{
		gNum:      number,
		taskTotal: total,
		tasks:     make(chan func() interface{}, total), // 缓冲通道个数是total,类型是func() interface{}
	}
}

// Start 开始执行任务
func (q *Queue) Start() {
	defer close(q.tasks) // 任务执行完毕后,关闭通道

	// 设置计数信号个数
	q.wg.Add(q.taskTotal)

	// 通过goroutineNumber个goroutine来执行task
	for i := 0; i < q.gNum; i++ {
		runtime.Gosched() // 让出cpu给其他goroutine
		go q.work()       // 对每个任务开启独立goroutine执行
	}

	// 等待goroutine执行完毕
	q.wg.Wait()

	// 当所有的任务执行完毕后回调
	if q.finishedCallback != nil {
		q.finishedCallback()
	}
}

// work 执行任务
func (q *Queue) work() {
	for {
		// 不断取出任务,直到chan关闭
		task, ok := <-q.tasks
		if !ok {
			break
		}

		res := task()
		// 完成一个task立即回调
		if q.taskCallback != nil {
			q.taskCallback(res)
		}

		q.wg.Done()
	}
}

// Add 添加任务
func (q *Queue) Add(task func() interface{}) {
	if len(q.tasks) <= q.taskTotal-1 { // 防止缓冲通道个数超出边界个数total
		q.tasks <- task
	}
}

// SetTaskCallback 设置单个任务执行后的回调函数
func (q *Queue) SetTaskCallback(callback func(res interface{})) {
	q.taskCallback = callback
}

// SetFinishedCallback 所有任务完成后,回调函数
func (q *Queue) SetFinishedCallback(callback func()) {
	q.finishedCallback = callback
}
