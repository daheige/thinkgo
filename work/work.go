/*
利用无缓冲chan创建goroutine池来控制一组task的执行
1. work 包的目的是展示如何使用无缓冲的通道来创建一个 goroutine 池,这些 goroutine 执行
并控制一组工作,让其并发执行。在这种情况下,使用无缓冲的通道要比随意指定一个缓冲区大
小的有缓冲的通道好,因为这个情况下既不需要一个工作队列,也不需要一组 goroutine 配合执行。
无缓冲的通道保证两个 goroutine 之间的数据交换。

2. 这种使用无缓冲的通道的方法允许使用者知道什么时候 goroutine 池正在执行工作,
而且如果池里的所有goroutine 都忙,无法接受新的工作的时候,也能及时通过通道来通知调用者。
使用无缓冲的通道不会有工作在队列里丢失或者卡住,所有工作都会被处理。
*/
package work

import (
	"log"
	"os"
	"sync"
)

// Worker worker必须满足Task方法
type Worker interface {
	Task()
}

// Pool提供一个goroutine池,可以完成任何已提交的worker任务
type Pool struct {
	work   chan Worker
	wg     sync.WaitGroup
	logger Logger
}

// Logger log interface
type Logger interface {
	Println(msg ...interface{})
}

var LogEntry Logger = log.New(os.Stderr, "", log.LstdFlags)

// New 创建一个工作池
func New(gNum int) *Pool {
	p := &Pool{
		work: make(chan Worker), // 无缓冲通道
	}

	p.wg.Add(gNum) // 最大goroutine个数
	for i := 0; i < gNum; i++ {
		// 开启独立的goroutine来执行任务
		go func(p *Pool) {
			defer catchRecover()

			defer p.wg.Done() // 执行完毕后计数信号量减去1

			// for...range会一直阻塞,直到从work通道中收到一个Worker接口值
			for w := range p.work {
				w.Task() // 执行任务
			}
		}(p)
	}

	return p
}

// Add 生产者:采用无缓冲通道提交任务到工作池
// 当任务提交后，消费者就会立即执行任务，p.wg计数器数量减去1
// w是一个接口值,必须是具体实现类型的一个实例指针
func (p *Pool) Add(w Worker) {
	p.work <- w
}

// Shutdown 等待所有的goroutine执行完毕,它关闭了 work 通道
// 这会导致所有池里的 goroutine 停止工作
// 调用pg.wg 的 Wait 方法,会等待所有 goroutine 终止
func (p *Pool) Shutdown() {
	close(p.work) // 关闭通道会让所有池里的goroutine全部停止
	p.wg.Wait()   // 等待所有的goroutine执行完毕
	LogEntry.Println("all goroutine task finish")
}

// catchRecover 捕获异常或者panic处理
func catchRecover() {
	if err := recover(); err != nil {
		LogEntry.Println("exec worker error: ", err)
	}
}
