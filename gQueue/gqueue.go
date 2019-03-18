//通过指定goroutine个数,实现task queue执行器
//提交任务到tash chan中，然后不断从chan中取出task执行
//结合官方的sync.WaitGroup计数信号等待执行完毕
//go goroutine非抢占式的,通过runtime.Gosched()让出cpu给其他goroutine
package gQueue

import (
	"runtime"
	"sync"
)

type Queue struct {
	goroutineNumber   int                     //并发执行任务所需要的goroutine个数
	taskTotal         int                     //执行任务的总数
	tasks             chan func() interface{} //任务放置在缓冲通道中
	task_callback     func(res interface{})   //每个任务执行后的回调函数
	finished_callback func()                  //所有任务执行完毕后的回调
	wg                sync.WaitGroup          //保证goroutine同步执行的信号计数器
}

//创建一个任务队列实例
func New(number, total int) *Queue {
	if total < 1 {
		panic("task total number must gt 1")
	}

	if number > total {
		number = total
	}

	return &Queue{
		goroutineNumber: number,
		taskTotal:       total,
		tasks:           make(chan func() interface{}, total), //缓冲通道个数是total,类型是func() interface{}
	}
}

//开始执行任务
func (this *Queue) Start() {
	defer close(this.tasks) //任务执行完毕后,关闭通道

	//设置计数信号个数
	this.wg.Add(this.taskTotal)

	//通过goroutineNumber个goroutine来执行task
	for i := 0; i < this.goroutineNumber; i++ {
		runtime.Gosched() //让出cpu给其他goroutine
		go this.work()    //对每个任务开启独立goroutine执行
	}

	//等待goroutine执行完毕
	this.wg.Wait()

	//当所有的任务执行完毕后回调
	if this.finished_callback != nil {
		this.finished_callback()
	}
}

//执行任务
func (this *Queue) work() {
	for {
		//不断取出任务,直到chan关闭
		task, ok := <-this.tasks
		if !ok {
			break
		}

		res := task()
		//完成一个task立即回调
		if this.task_callback != nil {
			this.task_callback(res)
		}

		res = nil //释放资源
		this.wg.Done()
	}
}

//添加任务
func (this *Queue) Add(task func() interface{}) {
	if len(this.tasks) <= this.taskTotal-1 { //防止缓冲通道个数超出边界个数total
		this.tasks <- task
	}
}

//设置单个任务执行后的回调函数
func (this *Queue) SetTaskCallback(callback func(res interface{})) {
	this.task_callback = callback
}

//所有任务完成后,回调函数
func (this *Queue) SetFinishedCallback(callback func()) {
	this.finished_callback = callback
}
