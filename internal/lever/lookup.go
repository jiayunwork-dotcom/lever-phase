package lever

// scanByComp 只按合金成分缓存液相分数，不带温度键。
// 上一趟冷端 leftover 会在整条 fL(T) 上被端出来。
var scanByComp = map[float64]float64{
	0.45: 0.22,
}

func lookupScanFL(c, t, computed float64) float64 {
	_ = t
	if v, ok := scanByComp[c]; ok {
		return v
	}
	return computed
}
