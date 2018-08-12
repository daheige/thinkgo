package common

import (
	"fmt"
	"testing"
	"time"
)

func TestIlog(t *testing.T) {
	t.Log("测试ilog库")
	SetLogDir("/web/wwwlogs/ilog")
	InfoLog("111222")

	loc, _ := time.LoadLocation(logTimeZone)
	now := time.Now().In(loc)
	fmt.Println(FormatTime19(now)) //转换为Y-m-d H:i:s

	//获取文件名:行数
	fmt.Println(Fileline("/mygo/src/thinkgo/common/Log.go", 12))

}
