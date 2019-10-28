package common

import "strconv"

//================str,int,int64,float64 conv func=======================
// IntToStr int-->string
func IntToStr(n int) string {
	return strconv.Itoa(n)
}

// StrToInt string-->int
func StrToInt(s string) int {
	if i, err := strconv.Atoi(s); err != nil {
		return 0
	} else {
		return i
	}
}

// Int64ToStr int64-->string
func Int64ToStr(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}

// StrToInt64 string--> int64
func StrToInt64(s string) int64 {
	if i64, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0
	} else {
		return i64
	}
}

// StrToFloat64 string--->float64
func StrToFloat64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	return f
}

// Float64ToStr float64 to string
// 'e' (-d.dddde±dd，十进制指数)
func Float64ToStr(f64 float64) string {
	return strconv.FormatFloat(f64, 'e', -1, 64)
}
