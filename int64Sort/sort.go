package int64Sort

import "sort"

// int64类型的切片排序
type Int64 []int64

func NewInt64Slice(s []int64) Int64 {
	return Int64(s)
}

func (s Int64) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s Int64) Len() int {
	return len(s)
}

func (s Int64) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Int64Sort int64 排序
func Int64Sort(s []int64) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
