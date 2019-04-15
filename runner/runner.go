/**runner用于按照顺序，执行程序任务操作，可作为cron作业或定时任务
runner 包可用于展示如何使用通道来监视程序的执行时间,如果程序运行时间太长,也可以
用 runner 包来终止程序。当开发需要调度后台处理任务的程序的时候,这种模式会很有用。这
个程序可能会作为 cron 作业执行,或者在基于定时任务的云环境(如 iron.io)里执行。
补充说明：
可能作为cron作业或基于定时任务,可以控制程序执行时间
使用通道来监控程序的执行时间，生命周期，甚至终止程序等。
我们这个程序叫runner，我们可以称之为执行者，它可以在后台执行任何任务
而且我们还可以控制这个执行者，比如强制终止它等
此外这个执行者也是一个很不错的模式，比如我们写好之后，交给定时任务去执行即可
比如cron，这个模式我们还可以扩展更高效率的并发，更多灵活的控制程序的生命周期
更高效的监控等。
*/
package runner

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrorTimeout = errors.New("received timeout")
var ErrInterrupt = errors.New("received interrupt signal")

//声明一个runner
type Runner struct {
	tasks      []func() error   //需要执行的任务func,依次执行的任务
	complete   chan error       //报告处理任务已完成
	timeout    <-chan time.Time //所有任务多久可以执行完毕，只能接受通道中的值
	interrupt  chan os.Signal   //可以控制强制终止的信号
	allErrors  map[int]error    //发生错误的task index对应的错误
	lastTaskId int              //最后一次完成的任务id
}

//定义一个工厂函数创建runner
func New(t time.Duration) *Runner {
	return &Runner{
		complete:  make(chan error, 1),     //有缓冲通道，存放所有任务运行后的结果状态
		timeout:   time.After(t),           //time.After返回time.Time类型的通道
		interrupt: make(chan os.Signal, 1), //声明一个中断信号
	}
}

func NewWithoutTime() *Runner {
	return &Runner{
		complete:  make(chan error, 1),
		interrupt: make(chan os.Signal, 1), //声明一个中断信号
	}
}

//将需要执行的任务添加到runner中
func (r *Runner) Add(tasks ...func() error) {
	r.tasks = append(r.tasks, tasks...) //r.tasks是一个函数切片类型
}

//获取已经完成的任务ids
func (r *Runner) ErrorTaskIds() map[int]error {
	return r.allErrors
}

//获取最后一次完成任务id
func (r *Runner) GetLastTaskId() int {
	return r.lastTaskId
}

//运行一个个任务,如果出错就返回错误信息
//run方法也很简单，会使用for循环，不停的运行任务，在运行的每个任务之前，都会检测是否收到了中断信号
//如果没有收到，则继续执行,一直到执行完毕，返回nil；
//如果收到了中断信号，则直接返回中断错误类型，任务执行终止
func (r *Runner) Run() error {
	for key, task := range r.tasks {
		if r.isInterrupt() {
			return ErrInterrupt
		}

		// log.Println("current run task id: ", key)
		e := task() //运行任务
		if e != nil {
			r.allErrors[key] = e
		}

		r.lastTaskId = key
	}

	return nil
}

//开始执行所有的任务
func (r *Runner) Start() error {
	//接受runner中断或者程序中断信号
	signal.Notify(r.interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	//让任务在协程中运行，当运行接受后，将运行结果给complete
	go func() {
		r.complete <- r.Run() //运行返回错误或者nil顺利完成
	}()

	// 判断哪个通道可操作
	select {
	case err := <-r.complete:
		return err
	case <-r.timeout: //超时，任务未执行完毕
		return ErrorTimeout
	}
}

//检查是否接受到操作系统的中断信号
//基于select的多路复用，select和switch很像，只不过它的每个case都是一个通信操作。
//那么到底选择哪个case块执行呢？原则就是哪个case的通信操作可以执行就执行哪个，如果同时有多个可以执行的case，那么就随机选择一个执行。
// 针对我们方法中，如果r.interrupt中接受不到值，就会执行default语句块，返回false
// 一旦r.interrupt中可以接收值，就会通知Go Runtime停止接收中断信号，然后返回true。
// 这里如果没有default的话，select是会阻塞的，直到r.interrupt可以接收值为止
func (r *Runner) isInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false //默认继续运行任务
	}
}
