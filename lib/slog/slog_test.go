package slog

import (
    "testing"
)

func TestLog(t *testing.T) {
    t.Log("开始测试")
    SetLogDir("/web/wwwlogs/slog")
    Info("fefefe", "test") //Info: 2018/06/02 22:05:20 fefefe
}
