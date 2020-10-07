package xsort

import "sort"

// Int64Slice int64类型的切片排序
type Int64Slice []int64

// NewInt64Slice 创建一个int64切片排序slice
func NewInt64Slice(s []int64) Int64Slice {
	return Int64Slice(s)
}

//Less less
func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

// Len return len
func (s Int64Slice) Len() int {
	return len(s)
}

// Swap swap
func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Int64QuickSort 对[]int64快速排序
func Int64QuickSort(s []int64) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// Int64StableSort 对[]int64稳定排序
// 从小到大，当元素相同时候，保持原有index顺序
func Int64StableSort(s []int64) {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

// ============对[]Type类型的切片进行排序=========================

// SliceStable 对切片类型的s进行稳定排序
// SliceStable sorts the provided slice given the provided less
// function while keeping the original order of equal elements.
// The function panics if the provided interface is not a slice.
func SliceStable(s interface{}, less func(i, j int) bool) {
	sort.SliceStable(s, less)
}

// Slice 对s切片类型的数据进行排序,当出现相同元素的话，采用快速排序
// The sort is not guaranteed to be stable.
// For a stable sort, use SliceStable.
func Slice(s interface{}, less func(i, j int) bool) {
	sort.Slice(s, less)
}
