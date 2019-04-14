package work

import (
	"log"
	"runtime"
	"sync"
)

//-----------------------worker interface----------------
// worker interface定义
// 只用任务对象上实现了Do方法就是一个Worker接口类型
type Worker interface {
	Do()
}

// WorkerPool定义
type WorkerPool struct {
	work chan Worker
	wg sync.WaitGroup
	gNum int
}

// New 创建workerPool实例
// 指定数量的goroutine来运行任务池中的任务
func New(gNum int) *WorkerPool {
	wp := &WorkerPool{
		work: make(chan Worker),
		wg: sync.WaitGroup{},
		gNum:gNum,
	}

	return wp
}

// Add添加任务worker到work通道中
func (wp WorkerPool) Add(w Worker) {
	wp.work <- w
}

//开始运行作业池中的任务
func (wp WorkerPool) run(){
	//独立协程运行task
	wp.wg.Add(wp.gNum) //最大goroutine个数
	for i := 0; i < wp.gNum; i++ {
		runtime.Gosched() //让出控制权给其他的goroutine,在逻辑上形成并发,只有出让cpu时，另外一个协程才会运行

		go func() {
			defer wp.Recover() //异常捕获处理
			defer wp.wg.Done() //当前goroutine执行完毕后，计数器减去1

			for task := range wp.work {
				task.Do()
			}
		}()
	}
}

// Wait方法调用,等待所有goroutine终止
func (wp WorkerPool) Wait(){
	wp.run() //开始运行任务
	wp.wg.Wait() //等待任务运行完毕
}

// goroutine recover捕获
func (wp WorkerPool) Recover(){
	if err := recover(); err != nil {
		log.Println("exec recover error: ",err)
	}
}