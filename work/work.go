//workerPool实现百万级的并发
//go程序开发过程中,通过简单的调用go func 函数来开启协程,容易导致程序死锁
//并且会无限制的开启groutine,groutine数量激增的情况下并发性能会明显下降
// 所以需要考虑使用工作池来控制协程数量,以达到高并发的效果.
package work

import (
	"log"
)

//定义任务接口,只要接口类型实现了Do()方法就实现了job接口
// 一个数据接口，所有的数据都要实现该接口，才能被传递进来
//实现Job接口的一个数据实例，需要实现一个Do()方法，对数据的处理就在这个Do()方法中
type Job interface {
	Do() error //如果Do方法没有错误返回，可以返回nil
}

//Job通道：
/*这里有两个Job通道：
	1、WorkerPool的Job channel，用于调用者把具体的数据写入到这里，WorkerPool读取。
	2、Worker的Job channel，当WorkerPool读取到Job，并拿到可用的Worker的时候，
      会将Job实例写入该Worker的Job channel，用来直接执行Do()方法。
*/

//-----------worker---------
//每一个被初始化的worker都会在后期单独占用一个协程
//初始化的时候会先把自己的JobQueue传递到Worker通道中，
//然后阻塞读取自己的JobQueue，读到一个Job就执行Job对象的Do()方法。
type Worker struct {
	JobQueue chan Job
}

// NewWorker初始化Worker
func NewWorker() Worker {
	return Worker{JobQueue: make(chan Job)}
}

//运行作业池中的任务
func (w Worker) Run(wq chan chan Job) {
	go func() {
		defer catchRecover()

		for {
			wq <- w.JobQueue //注册任务到wokerQueue中
			select {
			case job := <-w.JobQueue: //从通道中获取任务，执行任务
				if err := job.Do();err != nil{
					log.Println("exec task error: ",err.Error())
				}
			}
		}
	}()
}

//-------------------工作池(WorkerPool)------------
//初始化时会按照传入的num，启动num个后台协程，然后循环读取Job通道里面的数据，
//读到一个数据时，再获取一个可用的Worker，并将Job对象传递到该Worker的chan通道
/*工作池原理：
	1. 整个过程中 每个Worker都会被运行在一个协程中，在整个WorkerPool中就会有num可空闲的Worker
	2. 当来一条数据的时候，就会在工作池中去一个空闲的Worker去执行该Job，当工作池中没有可用的worker时
	就会阻塞等待一个空闲的worker。
 */
type WorkerPool struct {
	workerLen   int //WorkerPool中同时存在Worker的个数
	JobQueue    chan Job //WorkerPool的Job通道
	WorkerQueue chan chan Job
}

func NewWorkerPool(workerLen int) *WorkerPool {
	return &WorkerPool{
		workerLen:   workerLen, //作业个数
		JobQueue:    make(chan Job), //队列
		WorkerQueue: make(chan chan Job, workerLen),//作业队列,有缓冲取通道，大小为workerLen
	}
}

func (wp *WorkerPool) Run() {
	log.Println("init worker")

	//初始化worker
	for i := 0; i < wp.workerLen; i++ {
		worker := NewWorker()
		worker.Run(wp.WorkerQueue)
	}

	// 循环获取可用的worker,往空闲的worker中写入job
	go func() {
		defer catchRecover()

		for {
			select {
			case job := <-wp.JobQueue:
				worker := <-wp.WorkerQueue
				worker <- job //把job丢入workerQueue中
			}
		}
	}()
}

//捕获异常或者panic处理
func catchRecover(){
	if err := recover(); err != nil {
		log.Println("exec worker error: ",err)
	}
}