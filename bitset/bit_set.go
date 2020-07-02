/*
 * bitSet位图实现
 * 一、概述
	本文将讲述Bit-Map算法的相关原理,Bit-Map算法的一些利用场景
	例如BitMap解决海量数据寻找重复、判断个别元素是否在海量数据当中等问题
	最后说说BitMap的特点已经在各个场景的使用性。

二、Bit-Map算法
	先看看这样的一个场景（来自《编程珠玑》）：给一台普通PC，2G内存
	要求处理一个包含40亿个不重复并且没有排过序的无符号的int整数，
	给出一个整数，问如果快速地判断这个整数是否在文件40亿个数据当中？
	问题思考：
	40亿个int占（40亿*4）/1024/1024/1024 大概为14.9G左右，很明显内存只有2G
	放不下，因此不可能将这40亿数据放到内存中计算。要快速的解决这个问题最好的方案
	就是将数据搁内存了，所以现在的问题就在如何在2G内存空间以内存储着40亿整数。
	一个int整数在golang中是占4个字节的即要32bit位，如果能够用一个bit位来标识一个int整数
	那么存储空间将大大减少，算一下40亿个int需要的内存空间为40亿/8/1024/1024大概为476.83 mb
	这样的话我们完全可以将这40亿个int数放到内存中进行处理
*/

package bitset

import (
	"bytes"
	"fmt"
)

// IntSet An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
type IntSet struct {
	words []uint
}

const (
	bitNum = (32 << (^uint(0) >> 63))
)

// New new an entry
func New() *IntSet {
	return &IntSet{}
}

// Has reports whether the set contains the non-negative value x.
func (s *IntSet) Has(x int) bool {
	word, bit := x/bitNum, uint(x%bitNum)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
func (s *IntSet) Add(x int) {
	word, bit := x/bitNum, uint(x%bitNum)
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// UnionWith A与B的交集，合并A与B
// UnionWith sets s to the union of s and t.
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < bitNum; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}

				fmt.Fprintf(&buf, "%d", bitNum*i+j)
			}
		}
	}

	buf.WriteByte('}')
	return buf.String()
}

// Len len
func (s *IntSet) Len() int {
	var len int
	for _, word := range s.words {
		for j := 0; j < bitNum; j++ {
			if word&(1<<uint(j)) != 0 {
				len++
			}
		}
	}
	return len
}

// Remove 移除元素
func (s *IntSet) Remove(x int) {
	word, bit := x/bitNum, uint(x%bitNum)
	if s.Has(x) {
		s.words[word] ^= 1 << bit
	}
}

// Clear清空
func (s *IntSet) Clear() {
	s.words = []uint{}
}

// Copy copy value
func (s *IntSet) Copy() *IntSet {
	IntSet := &IntSet{
		words: []uint{},
	}

	IntSet.words = append(IntSet.words, s.words...)

	return IntSet
}

// AddAll 一次性添加多个int
func (s *IntSet) AddAll(args ...int) {
	for _, x := range args {
		s.Add(x)
	}
}

// IntersectWith A与B的并集，A与B中均出现
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i >= len(s.words) {
			continue
		}
		s.words[i] &= tword
	}
}

// DifferenceWith A与B的差集，元素出现在A未出现在B
func (s *IntSet) DifferenceWith(t *IntSet) {
	t1 := t.Copy() //为了不改变传参t，拷贝一份
	t1.IntersectWith(s)
	for i, tword := range t1.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		}
	}
}

// SymmetricDifference A与B的并差集，元素出现在A没有出现在B，或出现在B没有出现在A
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Elems 获取比特数组中的所有元素的slice集合
func (s *IntSet) Elems() []int {
	var elems []int
	for i, word := range s.words {
		for j := 0; j < bitNum; j++ {
			if word&(1<<uint(j)) != 0 {
				elems = append(elems, bitNum*i+j)
			}
		}
	}

	return elems
}
