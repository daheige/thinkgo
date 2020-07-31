// Package workpool for do task in work pool.
package workpool

import (
	"context"
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

	err := t.fn()
	if err != nil {
		logEntry.Println("exec task error: ", err)
	}
}

// Logger log record interface
type Logger interface {
	Println(args ...interface{})
}

// Pool task work pool
type Pool struct {
	// execInterval interval time after each task is executed
	// interval default 10ms
	execInterval   time.Duration
	entryChan      chan *Task     // task entry chan
	entryCap       int            // entry chan num
	jobChan        chan *Task     // job chan
	jobCap         int            // job chan num
	workerCap      int            // worker chan num
	logEntry       Logger         // logger interface
	stop           chan struct{}  // stop sem
	interrupt      chan os.Signal // interrupt signal
	entryCloseWait time.Duration  // close entry chan wait time,default 5s
	shutdownWait   time.Duration  // work pool shutdown wait time,default 3s
}

var (
	// defaultMaxEntryCap default max entry chan num.
	defaultMaxEntryCap = 10000

	// defaultMaxJobCap default max job chan cap.
	defaultMaxJobCap = 10000

	// defaultMaxWorker default max worker num
	// the number of workers depends on the specific business.
	defaultMaxWorker = 10000

	// defaultMinWorker default min worker.
	defaultMinWorker = 3

	// empty struct
	est = struct{}{}

	// dummy logger writes nothing.
	dummyLogger = LoggerFunc(func(...interface{}) {})
)

// LoggerFunc is a bridge between Logger and any third party logger.
type LoggerFunc func(...interface{})

// Println implements Logger interface.
func (f LoggerFunc) Println(args ...interface{}) { f(args...) }

// Option func Option to change pool.
type Option func(p *Pool)

// WithExecInterval interval time after each task is executed.
func WithExecInterval(t time.Duration) Option {
	return func(p *Pool) {
		p.execInterval = t
	}
}

// WithEntryNum task entry chan number.
func WithEntryCap(n int) Option {
	return func(p *Pool) {
		p.entryCap = n
	}
}

// WithJobCap job chan number.
func WithJobCap(n int) Option {
	return func(p *Pool) {
		p.jobCap = n
	}
}

// WithWorkerCap change worker num.
func WithWorkerCap(num int) Option {
	return func(p *Pool) {
		p.workerCap = num
	}
}

// WithLogger change logger entry.
func WithLogger(logEntry Logger) Option {
	return func(p *Pool) {
		p.logEntry = logEntry
	}
}

// WithEntryCloseWait close entry chan entryCloseWait time.
func WithEntryCloseWait(d time.Duration) Option {
	return func(p *Pool) {
		p.entryCloseWait = d
	}
}

// WithShutdownWait change shutdown entryCloseWait time.
func WithShutdownWait(d time.Duration) Option {
	return func(p *Pool) {
		p.shutdownWait = d
	}
}

// NewPool returns a pool.
func NewPool(opts ...Option) *Pool {
	p := &Pool{
		execInterval:   10 * time.Millisecond,
		workerCap:      defaultMinWorker,
		stop:           make(chan struct{}, 1),
		entryCloseWait: 5 * time.Second,
		shutdownWait:   3 * time.Second,
		interrupt:      make(chan os.Signal, 1),
		logEntry:       dummyLogger, // default logger entry.
	}

	// option functions.
	p.apply(opts...)

	if p.workerCap >= defaultMaxWorker {
		p.workerCap = defaultMaxWorker
	}

	if p.jobCap == 0 {
		// no buf for jobChan.
		p.jobChan = make(chan *Task)
	} else {
		if p.jobCap >= defaultMaxJobCap {
			p.jobCap = defaultMaxJobCap
		}

		p.jobChan = make(chan *Task, p.jobCap)
	}

	if p.entryCap == 0 {
		// no buf for entryChan.
		p.entryChan = make(chan *Task)
	} else {
		if p.entryCap >= defaultMaxEntryCap {
			p.entryCap = defaultMaxEntryCap
		}

		p.entryChan = make(chan *Task, p.entryCap)
	}

	return p
}

func (p *Pool) apply(opts ...Option) {
	for _, opt := range opts {
		opt(p)
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
		p.logEntry.Println("current worker id: ", id)

		// interval time after each task is executed.
		if p.execInterval > 0 {
			time.Sleep(p.execInterval)
		}
	}
}

// Run create workerCap goroutine to exec task.
func (p *Pool) Run() {
	p.logEntry.Println("exec task begin...")
	signal.Notify(p.interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	doneBuf := make(chan struct{}, p.workerCap)
	// create p.workerCap goroutine to do task
	for i := 0; i < p.workerCap; i++ {
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

		// Here you need to entryCloseWait for the task that has been sent to the entry chan
		// to ensure that it can be sent successfully
		ctx, cancel := context.WithTimeout(context.Background(), p.entryCloseWait)
		defer cancel()

		<-ctx.Done()

		close(p.entryChan)
	}()

	// entryCloseWait all job chan task to finish.
	for i := 0; i < p.workerCap; i++ {
		<-doneBuf
	}

	p.logEntry.Println("work pool shutdown success")
}

// Shutdown If all task are sent to the task entry chan, you can call this method to exit smoothly.
func (p *Pool) Shutdown() {
	// Create a deadline to manual exit entryCloseWait time.
	ctx, cancel := context.WithTimeout(context.Background(), p.entryCloseWait)
	defer cancel()

	// Doesn't block if no task, but will otherwise entryCloseWait
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
