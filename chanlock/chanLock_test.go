package chanlock

import (
	"log"
	"runtime"
	"sync"
	"testing"
)

var count = 1

func TestChanLock(t *testing.T) {
	log.Println("fefe")

	var wg sync.WaitGroup

	//抢占式的更新count，需要对count进行枷锁保护
	//如果不加锁，count每次执行后，值都不一样
	chLock := NewChanLock()

	nums := 1000
	wg.Add(nums) //建议一次性实现计数
	for i := 0; i < nums; i++ {
		runtime.Gosched() //让出当前cpu给其他goroutine执行

		go func() {
			defer wg.Done()
			chLock.Lock()
			defer chLock.Unlock()

			v := count
			log.Println("current count: ", v)
			v++
			count = v
		}()
	}

	log.Println("exec running....")
	wg.Wait()

	log.Println("count: ", count)
}

/**$ go test -v -test.run TestChanLock
2019/11/27 21:59:50 current count:  1000
2019/11/27 21:59:50 count:  1001
--- PASS: Test_chanLock (0.03s)
PASS
ok      github.com/daheige/thinkgo/chanlock     0.034s
*/
