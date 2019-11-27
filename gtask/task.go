package gtask

// do task with timeout or context cancel
// it will catch panic stack info when func exec panic
import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

//  TaskRes task返回的结果
type TaskRes struct {
	Err      error
	Result   chan interface{}
	CostTime float64
}

// DoTask 在独立携程中运行fn
// 这里返回结果设计为interface{},因为有时候返回结果可以是error
func DoTask(fn func() interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				log.Println("task exec panic error: ", err)
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn()
		res.Result <- r
	}()

	<-done

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithArgs 在独立携程中执行有参数的fn
func DoTaskWithArgs(fn func(args ...interface{}) interface{}, args ...interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				log.Println("task exec panic error: ", err)
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn(args...)
		res.Result <- r
	}()

	<-done

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithTimeout 采用done+select+time.After实现goroutine超时调用
func DoTaskWithTimeout(fn func() interface{}, timeout time.Duration) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn()
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-time.After(timeout):
		if res.Err == nil { //当执行过程中没有发生了panic的话，这里设置为任务超时错误
			res.Err = errors.New("task timeout")
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithContext 通过上下文context+done+select实现goroutine超时调用
func DoTaskWithContext(ctx context.Context, fn func() interface{}, timeout time.Duration) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	ctx2, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn()
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-ctx2.Done(): //超时了
		if res.Err == nil {
			res.Err = errors.New("task timeout")
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithTimeoutArgs 采用done+select+time.After实现goroutine超时调用
func DoTaskWithTimeoutArgs(fn func(args ...interface{}) interface{}, timeout time.Duration, args ...interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn(args...)
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-time.After(timeout):
		if res.Err == nil {
			res.Err = errors.New("task timeout")
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()
	return res
}

// DoTaskWithContextArgs 通过上下文context+done+select实现goroutine超时调用
func DoTaskWithContextArgs(ctx context.Context, fn func(args ...interface{}) interface{}, timeout time.Duration, args ...interface{}) *TaskRes {
	t := time.Now()
	done := make(chan struct{}, 1)
	res := &TaskRes{
		Result: make(chan interface{}, 1),
	}

	ctx2, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		defer func() {
			close(res.Result)
			close(done)
			if err := recover(); err != nil {
				res.Err = errors.New(fmt.Sprintf("%v", err))
			}
		}()

		r := fn(args...)
		res.Result <- r
	}()

	select {
	case <-done:
		log.Println("task has done")
	case <-ctx2.Done(): //超时了
		if res.Err == nil {
			res.Err = ctx2.Err()
		}
	}

	res.CostTime = time.Now().Sub(t).Seconds()

	return res
}
