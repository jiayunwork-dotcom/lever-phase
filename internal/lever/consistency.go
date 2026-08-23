package lever

import "math"

// Consistency 是对求解结果的一致性检查摘要，调用方与测试都可用。
type Consistency struct {
	TotalFraction float64
	AllInRange    bool
	SumToOne      bool
	HasLiquid     bool
	LiquidInRange bool
}

// CheckConsistency 校验结果满足：质量分数都在 [0,1]、总和为 1、
// 液相分数（如有）在 [0,1]。任何不满足都会在对应布尔位上反映。
func CheckConsistency(res PointResult) Consistency {
	c := Consistency{TotalFraction: res.TotalFraction()}
	c.AllInRange = true
	for _, ph := range res.Phases {
		if ph.MassFraction < 0 || ph.MassFraction > 1 {
			c.AllInRange = false
		}
	}
	c.SumToOne = math.Abs(c.TotalFraction-1) < 1e-9
	if len(res.Phases) == 0 {
		return c
	}
	liq, ok := res.FractionOf("liquid")
	if ok {
		c.HasLiquid = true
		c.LiquidInRange = liq >= 0 && liq <= 1
	}
	return c
}

// SumOfFractions 返回一组相的质量分数之和。
func SumOfFractions(phases []PhaseFraction) float64 {
	sum := 0.0
	for _, ph := range phases {
		sum += ph.MassFraction
	}
	return sum
}
