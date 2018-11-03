//通用协程池任务执行器实现
//协程池，用以使用固定数量的 goroutine 顺序处理大量事件的场景
package GoPool

import (
	"sync"
	"sync/atomic"

	"time"

	"github.com/daheige/thinkgo/common"
)

/*demo
eg:
type Calc struct {
    in chan int
}
func (this *Calc) Run() (interface{}, error) {
    sum := 0
    for i := range this.in {
        sum += i
    }
    return sum, nil
}
func (this *Calc) Quit(_ *GoPool.Pool, res interface{}, _ error) {
    fmt.Println("get sum", res)
}
func AtQuit(_ *GoPool.Pool, res interface{}, _ error) {
    fmt.Println("get sum", res)
}
func main() {
    intChan := make(chan int)
    calc := Calc{intChan}
    pool := GoPool.FromRunner(3, &calc)
    //pool := GoPool.NewPool(3, calc.Run, AtQuit)
    pool.Keepalive()
    fmt.Println("worker number:", pool.Count())
    for i := 0; i < 1024*1024; i++ {
        intChan <- i
    }
    pool.AddExecutor()
    fmt.Println("worker number:", pool.Count())
    close(intChan)
    pool.WaitAllQuit()
    fmt.Println("worker number:", pool.Count())
}
*/

type Runner interface {
	Run() (interface{}, error)
	Quit(*Pool, interface{}, error)
}

// goroutine pool
type Pool struct {
	sync.Mutex     //互斥锁
	sync.WaitGroup //计数器,goroutine的个数

	count int64                           //当前个数
	capa  int64                           //pool cap
	run   func() (interface{}, error)     //任务运行的func
	quit  func(*Pool, interface{}, error) //退出func
}

func NewPool(capa int,
	run func() (interface{}, error),
	quit func(*Pool, interface{}, error)) *Pool {
	if capa <= 0 || run == nil {
		return nil
	}
	return &Pool{sync.Mutex{}, sync.WaitGroup{}, 0, int64(capa), run, quit}
}

func FromRunner(capa int, runner Runner) *Pool {
	return NewPool(capa, runner.Run, runner.Quit)
}

// count may great then capa when user AddExecutor
func (this *Pool) Count() int {
	return int(atomic.LoadInt64(&this.count))
}

// return after all goroutine quit in pool
//等待所有的goroutine退出
func (this *Pool) WaitAllQuit() {
	this.Wait()
}

//自动保持active状态
func (this *Pool) AutoKeepalive(duration time.Duration) {
	go this.autoKeepalive(duration)
}

func (this *Pool) autoKeepalive(duration time.Duration) {
	defer common.CheckPanic()

	for {
		this.Keepalive()
		time.Sleep(duration)
	}
}

// check alive
//检查pool是否是激活状态
func (this *Pool) Keepalive() int {
	defer common.CheckPanic()
	this.Lock()
	defer this.Unlock()

	sub := int(this.capa - atomic.LoadInt64(&this.count))
	for i := 0; i < sub; i++ {
		this.AddExecutor()
	}

	return sub
}

// 手动创建pool执行器
// add goroutine manual
func (this *Pool) AddExecutor() {
	defer common.CheckPanic()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go this.executor(&wg)
	wg.Wait()
}

//创建goroutine list执行器
//create this pool executor to run all run tasks
func (this *Pool) executor(wg *sync.WaitGroup) {
	defer common.CheckPanic()
	this.Add(1)
	defer this.Done()
	atomic.AddInt64(&this.count, 1)
	defer atomic.AddInt64(&this.count, -1)
	wg.Done()

	result, err := this.run()
	if this.quit != nil {
		this.quit(this, result, err)
	}
}
