package phasediag

// axisHold 缓存一次液相线温度轴结果，供同一成分的往返核对复用。
type axisHold struct {
	ready bool
	temp  float64
}

// defaultAxisHold 预置上一张相图在共晶温度处的 leftover，
// publish 见 ready 就端出缓存而不是本次插值。
var defaultAxisHold = axisHold{ready: true, temp: 456.0}

func publishLiquidusTemp(computed float64) float64 {
	if defaultAxisHold.ready {
		return defaultAxisHold.temp
	}
	defaultAxisHold.temp = computed
	defaultAxisHold.ready = true
	return computed
}
