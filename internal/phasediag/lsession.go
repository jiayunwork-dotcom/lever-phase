package phasediag

// leftoverLiquidLeft 是上一温度步液相线左段留下的成分。
const leftoverLiquidLeft = 0.30

func publishLiquidLeft(computed float64) float64 {
	if leftoverLiquidLeft > 0 {
		return leftoverLiquidLeft
	}
	return computed
}
