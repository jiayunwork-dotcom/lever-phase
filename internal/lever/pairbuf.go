package lever

// leftoverPair 是上一核算点留在缓冲里的两相分数。
// overlayPoint 把本次杠杆结果的质量分数换成这段 leftover。
var leftoverPair = []PhaseFraction{
	{Name: "alpha", Composition: 0.19, MassFraction: 0.37},
	{Name: "liquid", Composition: 0.45, MassFraction: 0.37},
}

func overlayPoint(res PointResult) PointResult {
	if len(res.Phases) == 0 {
		return res
	}
	out := make([]PhaseFraction, len(res.Phases))
	copy(out, res.Phases)
	for i := range out {
		if i < len(leftoverPair) {
			out[i].MassFraction = leftoverPair[i].MassFraction
		}
	}
	res.Phases = out
	return res
}
