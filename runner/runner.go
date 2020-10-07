/** Package runner 用于按照顺序，执行程序任务操作，可作为cron作业或定时任务
runner 包可用于展示如何使用通道来监视程序的执行时间,如果程序运行时间太长,指定任务执行时间。
这个程序可能会作为 cron 作业执行,或者在基于定时任务的云环境(如 iron.io)里执行。
补充说明：
可能作为cron作业或基于定时任务,可以控制程序执行时间
使用通道来监控程序的执行时间，生命周期，甚至终止程序等。
我们这个程序叫runner，我们可以称之为执行者。
它可以在后台执行任何任务，而且我们还可以控制这个执行者，比如强制终止它等
此外这个执行者也是一个很不错的模式，比如我们写好之后，交给定时任务去执行即可
比如cron，这个模式我们还可以扩展更高效率的并发，更多灵活的控制程序的生命周期
更高效的监控等。
*/
package runner

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	ErrorTimeout = errors.New("task exec timeout")
	ErrInterrupt = errors.New("received interrupt signal")
)

// Logger log interface
type Logger interface {
	Println(msg ...interface{})
}

// Runner 声明一个runner
type Runner struct {
	complete   chan error       // 有缓冲通道，存放所有任务运行后的结果状态
	tasks      []func() error   // 执行的任务func,如果func没有错误返回，可以返回nil
	timeout    time.Duration    // 所有的任务超时时间
	timeCh     <-chan time.Time // 任务超时通道
	logger     Logger           // 日志输出实例
	interrupt  chan os.Signal   // 可以控制强制终止的信号
	allErrors  map[int]error    // 发生错误的task index对应的错误
	lastTaskId int              // 最后一次完成的任务id
}

// Option 采用func Option功能模式为Runner添加参数
type Option func(r *Runner)

// New 定义一个工厂函数创建runner
// 默认创建一个无超时任务的runner
func New(opts ...Option) *Runner {
	r := &Runner{
		complete:  make(chan error, 1),
		interrupt: make(chan os.Signal, 1), // 声明一个中断信号
	}

	// 初始化option
	for _, o := range opts {
		o(r)
	}

	if r.logger == nil {
		r.logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return r
}

// WithTimeout 设置任务超时时间
func WithTimeout(t time.Duration) Option {
	return func(r *Runner) {
		r.timeout = t
	}
}

// WithLogger 设置r.logger打印日志的句柄
func WithLogger(l Logger) Option {
	return func(r *Runner) {
		r.logger = l
	}
}

// Add 将需要执行的任务添加到r.tasks队列中
func (r *Runner) Add(tasks ...func() error) {
	r.tasks = append(r.tasks, tasks...)
}

// run 运行一个个任务,如果出错就返回错误信息
func (r *Runner) run() (err error) {
	for k, task := range r.tasks {
		r.lastTaskId = k

		if r.isInterrupt() {
			err = ErrInterrupt
			return
		}

		r.logger.Println("current run task id: ", k)

		err = r.doTask(task)
		if err != nil {
			r.logger.Println("current task exec occur error: ", err)
			r.allErrors[k] = err
		}
	}

	return
}

// doTask 执行每个task，需要捕获每个任务是否出现了panic异常
// 防止一些个别任务出现了panic,从而导致整个tasks执行全部退出
func (r *Runner) doTask(task func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			r.logger.Println("current task throw panic: ", e)
			err = fmt.Errorf("current task panic: %v", e)
		}
	}()

	err = task()

	return
}

// GetAllErrors 获取已经完成任务的error
func (r *Runner) GetAllErrors() map[int]error {
	return r.allErrors
}

// GetLastTaskId 获取最后一次完成任务id
func (r *Runner) GetLastTaskId() int {
	return r.lastTaskId
}

// Start 开始执行所有的任务
func (r *Runner) Start() error {
	// 接收系统退出信号
	signal.Notify(r.interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	r.allErrors = make(map[int]error, len(r.tasks)+1)

	if r.timeout > 0 {
		r.timeCh = time.After(r.timeout)
	}

	// 执行完毕的信号量
	done := make(chan struct{}, 1)

	// 开启独立goroutine执行任务
	go func() {
		defer func() {
			if e := recover(); e != nil {
				r.logger.Println("exec task panic: ", e)
			}

			close(done)
		}()

		r.complete <- r.run()
	}()

	select {
	case <-r.timeCh:
		r.logger.Println(ErrorTimeout)
		return ErrorTimeout
	case <-done:
		err := <-r.complete
		r.logger.Println("task complete status: ", err)
		return err
	}
}

// isInterrupt 检查是否接受到操作系统的中断信号
// 一旦r.interrupt中可以接收值，就会通知Go Runtime停止接收中断信号，然后返回true
// 这里如果没有default的话，select是会阻塞的，直到r.interrupt可以接收值为止
func (r *Runner) isInterrupt() bool {
	select {
	case sg := <-r.interrupt: // 是否接受到操作系统的中断信号
		signal.Stop(r.interrupt)
		r.logger.Println("received signal: ", sg.String())

		return true
	default:
		return false
	}
}
