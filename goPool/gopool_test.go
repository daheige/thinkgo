package goPool

import (
	"fmt"
	"testing"
)

type Calc struct {
	in chan int
}

func (this *Calc) Run() (interface{}, error) {
	sum := 0
	for i := range this.in {
		fmt.Printf("开始执行第%d次求和\n", i)
		sum += i
		fmt.Printf("当前%d累加的sum和: %d\n", i, sum)
	}

	// fmt.Println(sum)
	return sum, nil
}

//退出获取结果
func (this *Calc) Quit(_ *Pool, res interface{}, _ error) {
	fmt.Println("get sum", res)
}

func AtQuit(_ *Pool, res interface{}, _ error) {
	fmt.Println("get sum", res)
}

func Test_goroutine_pool(t *testing.T) {
	intChan := make(chan int)
	calc := Calc{intChan}

	//创建3个goroutine pool
	pool := FromRunner(3, &calc)
	//pool := NewPool(3, calc.Run, AtQuit)
	pool.Keepalive() //让3个goroutine保持激活状态
	fmt.Println("worker number:", pool.Count())

	//依次执行1048575次累加
	for i := 0; i < 1024*1024; i++ {
		intChan <- i
	}
	pool.AddExecutor() //创建pool执行器
	fmt.Println("worker number:", pool.Count())
	close(intChan)
	pool.WaitAllQuit() //等待所有的goroutine执行完毕

	fmt.Println("worker number:", pool.Count())
	t.Log("success")
}

/**
time go test -v
当前1048575累加的sum和: 184003456399
get sum 184003456399
worker number: 0
--- PASS: Test_goroutine_pool (17.18s)
	gopool_test.go:53: success
PASS
ok  	github.com/daheige/thinkgo/goPool	17.186s

real	0m17.554s
user	0m5.174s
sys	0m5.701s
*/
