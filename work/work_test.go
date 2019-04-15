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

//实现了Worker接口
func (this *myName) Task() {
	log.Println("current index", this.index)
	log.Println("your name is ", this.name)
}

type NamePrinter struct {
	name  string
	index int
}

//实现了worker接口
func (this *NamePrinter) Task() {
	log.Printf("执行第%d次", this.index)
	log.Println("current name:", this.name)
}

/**$ go test -v -test.run=TestPool
2019/04/15 21:41:21 all goroutine task finish
PASS
ok  	github.com/daheige/thinkgo/work	0.213s
*/
func TestPool(t *testing.T) {
	p := New(100) //使用10个goroutine 来创建工作池
	var wg sync.WaitGroup

	var taskNum = 100 * 100
	wg.Add(taskNum)
	for i := 0; i < taskNum; i++ {
		runtime.Gosched() //让出控制权给其他的goroutine
		u := myName{
			name:  "daheige",
			age:   28,
			index: i,
		}

		//异步的方式把任务提交到p中并运行
		// goroutine相互竞争,通过异步的方式将任务提交到p.pool中
		go func() {
			defer wg.Done()

			//一旦工作池里的 goroutine 接收到这个值
			// Add 方法就会返回。这也会导致 goroutine 将 WaitGroup 的计数递减,并终止 goroutine。
			p.Add(&u)
		}()
	}

	wg.Wait() //一旦wg计数器减到0了,就要执行p.Shutdown关闭所有的通道,关闭工作池

	// 让工作池停止工作,等待所有现有的工作完成
	p.Shutdown()
}

/**
$ go test -v -test.run=TestWorkPool
2019/04/15 21:42:02 all goroutine task finish
PASS
ok  	github.com/daheige/thinkgo/work	19.369s
*/
//添加任务后，执行Shutdown操作等待所有的Task都运行完毕后退出
func TestWorkPool(t *testing.T) {
	p := New(10) //开辟10个goroutine来运行任务

	num := 100 * 100 * 100
	for i := 0; i < num; i++ {
		u := myName{
			name:  "daheige",
			age:   num,
			index: i + 1,
		}

		p.Add(&u) //发送任务到goroutine池中，就会立即调用Task()执行
	}

	p.Shutdown()
}
