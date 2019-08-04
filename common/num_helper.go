package common

import (
	"math"
	"math/rand"
	"time"
)

//对浮点数进行四舍五入操作比如 12.125保留2位小数==>12.13
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

//生成m-n之间的随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	//随机种子
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}
