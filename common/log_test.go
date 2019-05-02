package common

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Log("测试ilog库")
	SetLogDir("/web/wwwlogs/ilog")

	var wg sync.WaitGroup
	var nums int = 30 //30w日志写入磁盘
	wg.Add(nums)      //一次性计数器设置，保证独立携程都成处理完毕

	for i := 0; i < nums; i++ {
		go func() {
			defer wg.Done()

			InfoLog("111222")
			DebugLog("this is debug: 111222")
			ErrorLog("error msg: 111222")
			NoticeLog("notice msg: 111222")
			WarnLog("warning: 111222")
			CritLog("crit msg: 111222")
			AlterLog("alter: 111222")
			EmergLog("emerg msg: 111222")
		}()
	}

	wg.Wait()
	log.Println("write log success")

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	fmt.Println(now.Format(logtmFmtTime)) //转换为Y-m-d H:i:s

	//获取文件名
	fmt.Println(filepath.Base("/mygo/src/thinkgo/common/Log.go"))

}

/**
 * 测试写入240w日志
 * go test -v -test.run=TestLog
//将log.go 加锁方式改成 NewMutexLock 测试报告
--- PASS: TestLog (65.81s)
    log_test.go:13: 测试ilog库
PASS
ok      github.com/daheige/thinkgo/common       66.017s

//将log.go 加锁方式改成 NewChanLock 测试报告
	=== RUN   TestLog
--- PASS: TestLog (69.58s)
    log_test.go:13: 测试ilog库
PASS
ok      github.com/daheige/thinkgo/common       69.695s
在golang底层，channel的实现是通过互斥锁和数组的方式实现的
而且还有一些其他字段的同步，因此sync实现的乐观锁的效率比用chan实现互斥锁更快
在大量的读写或者大量的i/o操作下，sync互斥锁实现的乐观锁，效率相对来说高一点
*/
