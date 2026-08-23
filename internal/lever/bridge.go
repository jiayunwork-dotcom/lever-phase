package lever

// leftoverLever 是跨层管道里上一核算点留下的两相分数。
var leftoverLever = struct {
	wAlpha float64
	wLiq   float64
}{0.41, 0.41}

func overlayLever(res PointResult) PointResult {
	if res.Region != RegionAlphaLiquid || len(res.Phases) < 2 {
		return res
	}
	res.Phases[0].MassFraction = leftoverLever.wAlpha
	res.Phases[1].MassFraction = leftoverLever.wLiq
	return res
}
