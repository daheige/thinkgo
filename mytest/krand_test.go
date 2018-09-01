package mytest

import (
	"testing"
	"thinkgo/common"
)

/**
$ go test -v -test.run TestKrand
=== RUN   TestKrand
--- PASS: TestKrand (0.00s)
    uuid_test.go:13: 纯数字: 933741
    uuid_test.go:14: 纯小写字母 uarmvg
    uuid_test.go:15: 纯大写字母 RIRFLG
    uuid_test.go:16: 数字大小写混合 ZkLF3FNRwcDk
PASS
ok      mytest  0.006s
*/
func TestKrand(t *testing.T) {
	s := common.Krand(6, 0)
	t.Log("纯数字:", s)
	t.Log("纯小写字母", common.Krand(6, 1))
	t.Log("纯大写字母", common.Krand(6, 2))
	t.Log("数字大小写混合", common.Krand(12, 3))
}
