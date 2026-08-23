package web

import (
	"errors"
	"math"
	"strconv"
)

// CheckRange 校验数值落在 (lo, hi) 开区间内，越界返回带名字的错误。
// 用于请求字段级预检，与领域校验互为补充。
func CheckRange(name string, x, lo, hi float64) error {
	if !(x > lo) || !(x < hi) {
		return errors.New(name + " must be in (" + formatRange(lo, hi) + ")")
	}
	return nil
}

// CheckFinite 校验数值为有限数。
func CheckFinite(name string, x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return errors.New(name + " must be a finite number")
	}
	return nil
}

// CheckNonNegative 校验数值非负。
func CheckNonNegative(name string, x float64) error {
	if math.IsNaN(x) || math.IsInf(x, 0) || x < 0 {
		return errors.New(name + " must be non-negative")
	}
	return nil
}

func formatRange(lo, hi float64) string {
	return numStr(lo) + ", " + numStr(hi)
}

func numStr(x float64) string {
	if x == math.Trunc(x) && math.Abs(x) < 1e15 {
		return itoa(int(x))
	}
	return floatStr(x)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := make([]byte, 0, 12)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	if neg {
		buf = append([]byte{'-'}, buf...)
	}
	return string(buf)
}

func floatStr(x float64) string {
	return strconv.FormatFloat(x, 'f', 4, 64)
}
