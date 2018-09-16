package work

import (
	"log"
	"runtime"
	"sync"
	"testing"
)

type myName struct {
	name  string
	age   int
	index int
}

func (this *myName) Task() {
	log.Println("current index", this.index)
	log.Println("your name is ", this.name)
}

/**
$ go test -v -test.run TestPool
--- PASS: TestPool2 (0.20s)
PASS
ok      thinkgo/work    0.379s
*/
func TestPool(t *testing.T) {
	p := New(10)

	for i := 0; i < 10000; i++ {
		u := myName{
			name:  "daheige",
			age:   28,
			index: i,
		}
		p.Run(&u) //发送任务到goroutine池中
	}

	p.Shutdown()
}

/**
$ go test -v -test.run TestPool2
--- PASS: TestPool2 (0.17s)
PASS
ok      thinkgo/work    0.176s
*/
//当我们加大wg计数器的个数为1000000,效率和TestPool测试的性能报告相差不大
func TestPool2(t *testing.T) {
	p := New(10) //开辟10个goroutine来运行任务
	var wg sync.WaitGroup
	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		runtime.Gosched() //让出控制权给其他的goroutine
		u := myName{
			name:  "daheige",
			age:   28,
			index: i + 1,
		}

		//将参数直接传递给go func中,减少闭包函数变量的查找路径
		//goroutine相互竞争,通过异步的方式将任务提交到p.pool中
		go func(p *Pool, wg *sync.WaitGroup) {
			p.Run(&u) //发送任务到goroutine池中
			wg.Done()
		}(p, &wg)
	}

	wg.Wait() //一旦wg计数器减到0了,就要执行p.Shutdown关闭所有的通道,关闭工作池
	p.Shutdown()
}
