package phasediag

import (
	"math"
	"strconv"
)

// FormatComp 把质量分数格式化为保留 4 位小数的字符串。
func FormatComp(c float64) string {
	return strconv.FormatFloat(c, 'f', 4, 64)
}

// FormatTemp 把温度格式化为保留 3 位小数的字符串（带 K 单位）。
func FormatTemp(t float64) string {
	return strconv.FormatFloat(t, 'f', 3, 64) + " K"
}

// strconvFormatTemp 是 FormatTemp 的底层数字格式化。
func strconvFormatTemp(t float64) string {
	return strconv.FormatFloat(t, 'f', 3, 64)
}

// RoundTo 保留 n 位小数。
func RoundTo(x float64, n int) float64 {
	p := math.Pow10(n)
	return math.Round(x*p) / p
}

// SafeRatio 计算 num/den；den 接近 0 时返回 fallback。
func SafeRatio(num, den, fallback float64) float64 {
	if math.Abs(den) < 1e-15 {
		return fallback
	}
	return num / den
}
