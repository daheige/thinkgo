package common

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Log("测试ilog库")
	SetLogDir("/web/wwwlogs/ilog")
	Info("111222")
	Debug("this is debug: 111222")
	Error("error msg: 111222")
	Notice("notice msg: 111222")
	Warn("warning: 111222")
	Crit("crit msg: 111222")
	Alter("alter: 111222")
	Emerg("emerg msg: 111222")

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	fmt.Println(now.Format(logtmFmtTime)) //转换为Y-m-d H:i:s

	//获取文件名
	fmt.Println(filepath.Base("/mygo/src/thinkgo/common/Log.go"))

}

/**
 * go test -v -test.run=TestLog
=== RUN   TestLog
2018-10-28
Log.go
--- PASS: TestLog (0.00s)
	log_test.go:11: 测试ilog库
PASS
ok  	thinkgo/common	0.007s
 * 日志格式：
2018-10-27 22:34:39 info log_test.go line:[13]:111222
2018-10-27 22:34:39 debug log_test.go line:[14]:this is debug: 111222
2018-10-27 22:34:39 error log_test.go line:[15]:error msg: 111222
2018-10-27 22:34:39 notice log_test.go line:[16]:notice msg: 111222
2018-10-27 22:34:39 warn log_test.go line:[17]:warning: 111222
2018-10-27 22:34:39 crit log_test.go line:[18]:crit msg: 111222
2018-10-27 22:34:39 alter log_test.go line:[19]:alter: 111222
2018-10-27 22:34:39 emerg log_test.go line:[20]:emerg msg: 111222
*/
