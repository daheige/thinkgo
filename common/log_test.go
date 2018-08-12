package common

import (
	"testing"
)

func TestIlog(t *testing.T) {
	t.Log("测试ilog库")
	SetLogDir("/web/wwwlogs/ilog")
	InfoLog("111222")

}
