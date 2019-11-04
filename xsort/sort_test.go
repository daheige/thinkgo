package xsort

import (
	"log"
	"sort"
	"testing"
)

type people struct {
	Name string
	Age  int
}

func TestSort(t *testing.T) {
	arr := NewInt64Slice([]int64{
		1, 2, 3, 45, 23, 11, 9, 8, 23, 12,
	})

	sort.Sort(arr)
	log.Println(arr)

	//对int64类型的切片,实现快速排序
	s := []int64{1, 2, 3, 2, 12, 12, 10, 9, 1}
	Int64QuickSort(s)

	log.Println(s)

	s2 := []int64{
		1, 2, 3, 2, 34, 12, 123, 9, 90, 9, 10, 11,
	}

	Int64StableSort(s2)

	log.Println(s2)

	p := []people{
		{"Alice", 25},
		{"Elizabeth", 75},
		{"Alice", 75},
		{"Bob", 75},
		{"Alice", 75},
		{"Bob", 25},
		{"Colin", 25},
		{"Elizabeth", 25},
	}

	// Sort by age preserving name order
	Slice(p, func(i, j int) bool {
		return p[i].Age < p[j].Age
	})

	log.Println(p)

	Slice(p, func(i, j int) bool {
		return p[i].Name < p[j].Name
	})

	log.Println("sort by name of p: ", p)

	p2 := []people{
		{"Alice", 25},
		{"Elizabeth", 75},
		{"Alice", 75},
		{"Bob", 75},
		{"Alice", 75},
		{"Bob", 25},
		{"Colin", 25},
		{"Elizabeth", 25},
	}

	// SliceStable(p2, func(i, j int) bool {
	// 	return p2[i].Name < p2[j].Name
	// })

	SliceStable(p2, func(i, j int) bool {
		return p2[i].Age < p2[j].Age
	})

	log.Println("stable sort of p2: ", p2)

}

/**
$ go test -v
=== RUN   TestSort
2019/11/04 23:07:29 [1 2 3 8 9 11 12 23 23 45]
2019/11/04 23:07:29 [1 1 2 2 3 9 10 12 12]
2019/11/04 23:07:29 [1 2 2 3 9 9 10 11 12 34 90 123]
2019/11/04 23:07:29 [{Alice 25} {Elizabeth 25} {Bob 25} {Colin 25} {Alice 75} {Bob 75} {Alice 75} {Elizabeth 75}]
2019/11/04 23:07:29 sort by name of p:  [{Alice 25} {Alice 75} {Alice 75} {Bob 25} {Bob 75} {Colin 25} {Elizabeth 25} {Elizabeth 75}]
2019/11/04 23:07:29 stable sort of p2:  [{Alice 25} {Bob 25} {Colin 25} {Elizabeth 25} {Elizabeth 75} {Alice 75} {Bob 75} {Alice 75}]
--- PASS: TestSort (0.00s)
PASS
ok  	github.com/daheige/thinkgo/xsort	0.003s
*/
