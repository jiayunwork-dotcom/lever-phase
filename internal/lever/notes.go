package lever

import (
	"strconv"
)

// shortFmt 把质量分数格式化为 4 位小数文本。
func shortFmt(w float64) string {
	return strconv.FormatFloat(w, 'f', 4, 64)
}

// ExplainRegion 返回相区的人类可读说明（英文），供页面与日志使用。
func ExplainRegion(r Region) string {
	switch r {
	case RegionLiquid:
		return "liquid: single liquid phase, w=1, composition equals alloy"
	case RegionAlpha:
		return "alpha: single alpha solid solution, w=1, composition equals alloy"
	case RegionBeta:
		return "beta: single beta solid solution, w=1, composition equals alloy"
	case RegionAlphaLiquid:
		return "alpha + liquid: lever law with alpha solidus and liquidus-left boundary at current T"
	case RegionBetaLiquid:
		return "beta + liquid: lever law with beta solidus and liquidus-right boundary at current T"
	case RegionAlphaBeta:
		return "alpha + beta: eutectic solid below TE, lever law with TE boundary compositions"
	default:
		return string(r)
	}
}

// FormatTemperature 输出温度文本。
func FormatTemperature(t float64) string {
	return strconv.FormatFloat(t, 'f', 3, 64)
}

// FormatComposition 输出成分文本。
func FormatComposition(c float64) string {
	return strconv.FormatFloat(c, 'f', 4, 64)
}
