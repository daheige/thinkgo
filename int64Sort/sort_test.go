package int64Sort

import (
	"log"
	"sort"
	"testing"
)

func TestInt64(t *testing.T) {
	arr := NewInt64Slice([]int64{
		1, 2, 3, 45, 23, 11, 9, 8, 23, 12,
	})

	sort.Sort(arr)
	log.Println(arr)

	//对int64类型的切片排序
	s := []int64{1, 2, 3, 2, 12, 12, 10, 9, 1}
	Int64Sort(s)

	log.Println(s)
}

/**
$ go test -v
=== RUN   TestInt64
2019/11/04 22:25:02 [1 2 3 8 9 11 12 23 23 45]
2019/11/04 22:25:02 [1 1 2 2 3 9 10 12 12]
--- PASS: TestInt64 (0.00s)
PASS
ok  	github.com/daheige/thinkgo/int64Sort	0.002s
*/
