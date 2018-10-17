/*
利用无缓冲chan创建goroutine池来控制一组task的执行
work 包的目的是展示如何使用无缓冲的通道来创建一个 goroutine 池,这些 goroutine 执行
并控制一组工作,让其并发执行。在这种情况下,使用无缓冲的通道要比随意指定一个缓冲区大
小的有缓冲的通道好,因为这个情况下既不需要一个工作队列,也不需要一组 goroutine 配合执行。
无缓冲的通道保证两个 goroutine 之间的数据交换。
这种使用无缓冲的通道的方法允许使用者知道什么时候 goroutine 池正在执行工作,而且如果池里的所有
goroutine 都忙,无法接受新的工作的时候,也能及时通过通道来通知调用者。
使用无缓冲的通道不会有工作在队列里丢失或者卡住,所有工作都会被处理。
*/
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
			p.wg.Done() //当p.work关闭后，会让所有的goroutine全部停止工作，计数器减去1
		}(p)
	}

	return p
}

//生产者:采用无缓冲通道提交任务到工作池
//当任务提交后，消费者就会立即执行任务，p.wg计数器数量减去1
//w是一个接口值,必须是具体实现类型的一个实例指针
func (p *Pool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有的goroutine执行完毕,它关闭了 work 通道,这会导致所有池里的 goroutine 停止工作
// 并调用 WaitGroup 的 Done 方法使得计数器减去1;
// 然后调用pg.wg 的 Wait 方法,这会让 Shutdown 方法等待所有 goroutine 终止
func (p *Pool) Shutdown() {
	close(p.work) //关闭通道会让所有池里的goroutine全部停止
	// log.Println("all goroutine task finish")
	p.wg.Wait() //等待所有的goroutine执行完毕
}
