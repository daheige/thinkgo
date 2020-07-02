// Package workpool for do task in work pool.
package workpool

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Task task struct.
type Task struct {
	fn func() error
}

// NewTask returns task,create a task entry.
func NewTask(fn func() error) *Task {
	return &Task{
		fn: fn,
	}
}

// run exec a task.
func (t *Task) run(logEntry Logger) {
	if t == nil {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			logEntry.Println("exec current task panic: ", e)
		}
	}()

	t.fn()
}

// Logger log record interface
type Logger interface {
	Println(args ...interface{})
}

// Pool task work pool
type Pool struct {
	entryChan    chan *Task     // task entry chan
	jobChan      chan *Task     // job chan
	workerNum    int            // worker num
	logEntry     Logger         // logger interface
	stop         chan struct{}  // stop sem
	interrupt    chan os.Signal // interrupt signal
	wait         time.Duration  // close entry chan wait time,default 5s
	shutdownWait time.Duration  // work pool shutdown wait time,default 3s
}

var (
	// defaultMaxEntryCap default max entry chan num.
	defaultMaxEntryCap = 10000

	// defaultMaxJobCap default max job chan cap.
	defaultMaxJobCap = 10000

	// defaultMaxWorker default max worker num
	// the number of workers depends on the specific business.
	defaultMaxWorker = 100

	est = struct{}{}

	// LogDebug log task run in which worker.
	LogDebug = false
)

// Option func Option to change pool
type Option func(p *Pool)

// NewPool returns a pool.
func NewPool(num int, ec int, jc ...int) *Pool {
	if num < 0 || num >= defaultMaxWorker {
		num = defaultMaxWorker
	}

	if ec < 0 || ec >= defaultMaxEntryCap {
		ec = defaultMaxEntryCap
	}

	// jobChan buf num
	var jobChanCap int
	if len(jc) > 0 && jc[0] > 0 {
		jobChanCap = jc[0]
	}

	if jobChanCap < 0 || jobChanCap >= defaultMaxJobCap {
		jobChanCap = defaultMaxJobCap
	}

	p := &Pool{
		workerNum:    num,
		stop:         make(chan struct{}, 1),
		wait:         5 * time.Second,
		shutdownWait: 3 * time.Second,
		interrupt:    make(chan os.Signal, 1),
	}

	if jobChanCap == 0 {
		// no buf for jobChan.
		p.jobChan = make(chan *Task)
	} else {
		p.jobChan = make(chan *Task, ec)
	}

	if ec == 0 {
		// no buf for entryChan.
		p.entryChan = make(chan *Task)
	} else {
		p.entryChan = make(chan *Task, ec)
	}

	if p.logEntry == nil {
		p.logEntry = log.New(os.Stderr, "", log.LstdFlags)
	}

	return p
}

// WithLogger change logger entry.
func WithLogger(logEntry Logger) Option {
	return func(p *Pool) {
		p.logEntry = logEntry
	}
}

// WithWaitTime close entry chan wait time.
func WithWaitTime(d time.Duration) Option {
	return func(p *Pool) {
		p.wait = d
	}
}

// WithWorkerNum change worker num.
func WithWorkerNum(num int) Option {
	return func(p *Pool) {
		p.workerNum = num
	}
}

// WithShutdownWait change shutdown wait time.
func WithShutdownWait(d time.Duration) Option {
	return func(p *Pool) {
		p.shutdownWait = d
	}
}

// AddTask add a task to p.entryChan.
func (p *Pool) AddTask(t *Task) {
	if t == nil {
		return
	}

	defer p.recovery()

	select {
	case <-p.stop:
		return
	default:
		p.entryChan <- t
	}
}

// BatchAddTask batch add task to p.entryChan.
func (p *Pool) BatchAddTask(t []*Task) {
	defer p.recovery()

	for k := range t {
		if t[k] == nil {
			continue
		}

		select {
		case <-p.stop:
			return
		default:
			p.entryChan <- t[k]
		}
	}
}

// exec exec task from job chan.
func (p *Pool) exec(id int, done chan struct{}) {
	defer p.recovery()

	defer func() {
		done <- est
		p.logEntry.Println("current worker id: ", id, "will exit...")
	}()

	// get task from JobChan to run.
	for task := range p.jobChan {
		task.run(p.logEntry)
		if LogDebug {
			p.logEntry.Println("current worker id: ", id)
		}
	}
}

// Run create workerNum goroutine to exec task.
func (p *Pool) Run() {
	p.logEntry.Println("exec task begin...")
	signal.Notify(p.interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	doneBuf := make(chan struct{}, p.workerNum)

	// create p.workerNum goroutine to do task
	for i := 0; i < p.workerNum; i++ {
		go p.exec(i+1, doneBuf)
	}

	// throw entry chan task to JobChan
	go func() {
		defer close(p.jobChan)

		// If the entry channel is closed, the block will be lifted until consumption is completed.
		for task := range p.entryChan {
			p.jobChan <- task
		}
	}()

	// listen interrupt signal for work pool graceful exit.
	go func() {
		// Block until we receive stop signal.
		sig := <-p.interrupt
		p.logEntry.Println("recv signal: ", sig.String())
		close(p.stop)

		// Here you need to wait for the task that has been sent to the entry chan
		// to ensure that it can be sent successfully
		ctx, cancel := context.WithTimeout(context.Background(), p.wait)
		defer cancel()

		<-ctx.Done()

		close(p.entryChan)
	}()

	// wait all job chan task to finish.
	for i := 0; i < p.workerNum; i++ {
		<-doneBuf
	}

	p.logEntry.Println("work pool shutdown success")
}

// Shutdown If all task are sent to the task entry chan, you can call this method to exit smoothly.
func (p *Pool) Shutdown() {
	// Create a deadline to manual exit wait time.
	ctx, cancel := context.WithTimeout(context.Background(), p.wait)
	defer cancel()

	// Doesn't block if no task, but will otherwise wait
	// until the timeout deadline.
	<-ctx.Done()

	p.interrupt <- syscall.SIGTERM
	p.logEntry.Println("work pool will shutdown...")
}

// recovery catch a recover.
func (p *Pool) recovery() {
	defer func() {
		if e := recover(); e != nil {
			p.logEntry.Println("exec panic: ", e)
		}
	}()
}
