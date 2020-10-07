package gnum

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Round 对浮点数进行四舍五入操作比如 12.125保留2位小数==>12.13
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

// RandInt64 生成m-n之间的随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	// 随机种子
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// Abs abs(num)返回绝对值，返回结果是float64
func Abs(number float64) float64 {
	return math.Abs(number)
}

// Rand 产生[m,n]区间的int随机数
// Range: [0, 2147483647]
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}

	// PHP: getrandmax()
	if int31 := 1<<31 - 1; max > int31 {
		panic("max: max can not be greater than " + strconv.Itoa(int31))
	}

	if min == max {
		return min
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max+1-min) + min
}

// Floor 产生一个向下取整的数字，对应php floor(num)
func Floor(value float64) float64 {
	return math.Floor(value)
}

// Ceil 产生一个向上取整的数字，php ceil(num)
func Ceil(value float64) float64 {
	return math.Ceil(value)
}

// Max 获得数字的最大值 max(1,2,3)
func Max(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}

	max := nums[0]
	for i := 1; i < len(nums); i++ {
		max = math.Max(max, nums[i])
	}

	return max
}

// Min min(1,2,3)返回最小值
func Min(nums ...float64) float64 {
	if len(nums) < 2 {
		panic("nums: the nums length is less than 2")
	}
	min := nums[0]
	for i := 1; i < len(nums); i++ {
		min = math.Min(min, nums[i])
	}
	return min
}

// IsNumeric 判断v是否是一个数字类型，php is_numeric(num)
// Numeric strings consist of optional sign, any number of digits,
// optional decimal part and optional exponential part.
// Thus +0123.45e6 is a valid numeric value.
// In PHP hexadecimal (e.g. 0xf4c3b00c) is not supported,
// but IsNumeric is supported.
func IsNumeric(v interface{}) bool {
	switch val := v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64, complex64, complex128:
		return true
	case string:
		if val == "" {
			return false
		}

		// Trim any whitespace
		val = strings.TrimSpace(val)
		if val[0] == '-' || val[0] == '+' {
			if len(val) == 1 {
				return false
			}
			val = val[1:]
		}

		// hex
		if len(val) > 2 && val[0] == '0' && (val[1] == 'x' || val[1] == 'X') {
			for _, h := range val[2:] {
				if !((h >= '0' && h <= '9') || (h >= 'a' && h <= 'f') || (h >= 'A' && h <= 'F')) {
					return false
				}
			}
			return true
		}

		// 0-9, Point, Scientific
		p, s, l := 0, 0, len(val)
		for i, v1 := range val {
			if v1 == '.' { // Point
				if p > 0 || s > 0 || i+1 == l {
					return false
				}
				p = i
			} else if v1 == 'e' || v1 == 'E' { // Scientific
				if i == 0 || s > 0 || i+1 == l {
					return false
				}

				s = i
			} else if v1 < '0' || v1 > '9' {
				return false
			}
		}

		return true
	}

	return false
}
