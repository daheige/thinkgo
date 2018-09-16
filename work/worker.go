//利用无缓冲chan创建goroutine池来控制一组task的执行
//无缓冲通道保证两个goroutine之间的数据交换
//同步执行一组动作
package work

import (
	"runtime"
	"sync"
)

//worker必须满足Tash方法
type Worker interface {
	Task()
}

//Pool提供一个goroutine池,可以完成任何已提交的Woker任务
type Pool struct {
	work chan Worker
	wg   sync.WaitGroup
}

//创建一个工作池
func New(maxGoroutines int) *Pool {
	p := &Pool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines) //最大goroutine个数
	for i := 0; i < maxGoroutines; i++ {
		runtime.Gosched() //让出控制权给其他的goroutine,在逻辑上形成并发,只有出让cpu时，另外一个协程才会运行

		//开启独立的goroutine来执行任务
		go func(p *Pool) {
			for w := range p.work { //for...range会一直阻塞,知道从work通道中收到一个Worker接口值
				w.Task() //执行任务
			}

			//执行完毕后计数信号量减去1
			p.wg.Done()
		}(p)
	}

	return p
}

//生产者:run提交工作到工作池
//w是一个接口值,必须是具体实现类型的一个实例指针
func (p *Pool) Run(w Worker) {
	p.work <- w
}

//等待所有的goroutine执行完毕
func (p *Pool) Shutdown() {
	close(p.work) //关闭通道会让所有池里的goroutine全部停止
	// log.Println("all goroutine task finish")
	p.wg.Wait() //等待所有的goroutine执行完毕
}
