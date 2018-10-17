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

func TestWork(t *testing.T) {
	var names = []string{
		"daheige",
		"phper",
		"golang",
	}

	p := New(2) //开启2个goroutine来处理任务

	var wg sync.WaitGroup
	var nums = 300            //执行的组数
	wg.Add(nums * len(names)) //开启10 * len(names)个goroutine
	for i := 0; i < nums; i++ {
		//迭代names 每一组都执行一次打印操作
		for _, name := range names {
			//go 是非抢占的，只有出让cpu时，另外一个协程才会运行
			runtime.Gosched() //让出控制权给其他的goroutine,让goroutine相互竞争
			np := NamePrinter{
				name:  name,
				index: i + 1,
			}

			//提交任务到p.work中
			//发送任务到p.work
			go func(p *Pool, wg *sync.WaitGroup) {
				p.Run(&np)
				wg.Done()
			}(p, &wg)
		}
	}

	wg.Wait() //一旦wg计数器减到0了,就要执行p.Shutdown关闭所有的通道,关闭工作池

	//让工作池停止工作,等待所有的工作完成
	p.Shutdown()
}

/**
$ time go test -v -test.run TestWork
--- PASS: TestWork (0.03s)
PASS
ok  	thinkgo/work	0.028s

real	0m0.340s
user	0m0.395s
sys	0m0.082s

*/

/**
$ go test -v -test.run TestPool
--- PASS: TestPool2 (0.20s)
PASS
ok      thinkgo/work    0.379s
*/
func TestPool(t *testing.T) {
	p := New(10) //使用10个goroutine 来创建工作池
	var wg sync.WaitGroup

	var taskNum = 10000
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
			//一旦工作池里的 goroutine 接收到这个值
			// Run 方法就会返回。这也会导致 goroutine 将 WaitGroup 的计数递减,并终止 goroutine。
			p.Run(&u)
			wg.Done()
		}()
	}
	wg.Wait() //一旦wg计数器减到0了,就要执行p.Shutdown关闭所有的通道,关闭工作池

	// 让工作池停止工作,等待所有现有的工作完成
	p.Shutdown()
}

/**
$ go test -v -test.run TestPool2
--- PASS: TestPool2 (0.17s)
PASS
ok      thinkgo/work    1.808s
*/
//无缓冲通道的同步执行任务
func TestPool2(t *testing.T) {
	p := New(10) //开辟10个goroutine来运行任务
	for i := 0; i < 100000; i++ {
		u := myName{
			name:  "daheige",
			age:   28,
			index: i + 1,
		}

		p.Run(&u) //发送任务到goroutine池中，就会立即调用Task()执行
	}

	p.Shutdown()
}
