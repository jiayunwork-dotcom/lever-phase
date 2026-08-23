package web

import (
	"math"
	"strconv"
)

// round4 保留 4 位小数，用于把内核数值收敛到展示精度。
func round4(x float64) float64 {
	return math.Round(x*1e4) / 1e4
}

// round3 保留 3 位小数。
func round3(x float64) float64 {
	return math.Round(x*1e3) / 1e3
}

// format4 输出 4 位小数文本。
func format4(x float64) string {
	return strconv.FormatFloat(round4(x), 'f', 4, 64)
}

// format3 输出 3 位小数文本。
func format3(x float64) string {
	return strconv.FormatFloat(round3(x), 'f', 3, 64)
}

// clip01 把数值夹紧到 [0,1]。
func clip01(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
