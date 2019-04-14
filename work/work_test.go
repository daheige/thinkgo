package work

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"testing"
	"time"
)

type Score struct {
	Num int
}

func (s *Score) Do() error {
	fmt.Println("num:", s.Num)
	time.Sleep(1 * 1 * time.Second)
	return nil
}

func TestWork(t *testing.T) {
	num := 10000 //指定goroutine个数
	// debug.SetMaxThreads(num + 1000) //设置最大线程数
	// 注册工作池，传入任务参数1 worker并发个数
	p := NewWorkerPool(num)
	p.Run() //运行作业池

	//注册任务到jobQueue中
	datanum := 100 * 100 * 2
	go func() {
		for i := 1; i <= datanum; i++ {
			sc := &Score{Num: i}
			p.JobQueue <- sc
		}
	}()

	for {
		fmt.Println("runtime.NumGoroutine() :", runtime.NumGoroutine()) //10004
		time.Sleep(1 * time.Second)
	}

	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	//window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	sig:= <-ch
	log.Println("signal: ",sig.String())
	log.Println("workerPool will exit...")
}
