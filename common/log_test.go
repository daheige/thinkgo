package common

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	t.Log("开始测试")
	loc := GetLoc("PRC")
	t.Log(time.Now().In(loc).Format(tmFmtMissMS))
	t.Log(time.Now().In(loc).Format("2006-01-02-15-04-05"))
	SetLogDir("/web/wwwlogs/golang")
	InfoLog("123456") //2018-06-02 21:44:27 log_test:14 info [123456]
}
