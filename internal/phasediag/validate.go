package phasediag

import (
	"fmt"
	"math"
)

// ValidateComposition 校验成分为有限数且落在 [0, 1] 质量分数范围内。
// 质量分数超出 0～1 属于非物理输入，必须报错而不是静默截断。
func ValidateComposition(c float64) error {
	if math.IsNaN(c) || math.IsInf(c, 0) {
		return fmt.Errorf("composition must be a finite number, got %v", c)
	}
	if c < 0 || c > 1 {
		return fmt.Errorf("composition c = %v is out of range [0, 1]", c)
	}
	return nil
}

// ValidateTemperature 校验温度为有限数且不低于绝对零度（0 K）。
// 温度低于 0 K 属于非物理输入，必须报错。
func ValidateTemperature(t float64) error {
	if math.IsNaN(t) || math.IsInf(t, 0) {
		return fmt.Errorf("temperature must be a finite number, got %v", t)
	}
	if t < 0 {
		return fmt.Errorf("temperature t = %v K is below absolute zero", t)
	}
	return nil
}

// ValidatePolylineNodes 校验一组结点能否构成合法相界折线：
// 成分必须严格递增，温度方向翻转至多一次，温度为有限非负值。
func ValidatePolylineNodes(nodes []Node) error {
	if len(nodes) < 2 {
		return fmt.Errorf("polyline needs at least two nodes, got %d", len(nodes))
	}
	p := Polyline{nodes: append([]Node(nil), nodes...)}
	return p.Validate()
}

// ValidateDiagramTopology 校验整张相图的拓扑一致性：共晶温度低于
// 两个端元熔点、共晶成分位于两端元溶解度之间、三条折线温度轴对齐。
func ValidateDiagramTopology(d PhaseDiagram) error {
	if d.EutecticT <= 0 {
		return fmt.Errorf("eutectic temperature %v K must be positive", d.EutecticT)
	}
	if d.EutecticC <= d.AlphaMax || d.EutecticC >= d.BetaMin {
		return fmt.Errorf("eutectic composition %v must lie strictly between alpha max %v and beta min %v", d.EutecticC, d.AlphaMax, d.BetaMin)
	}
	if d.EutecticT >= d.Liquidus.MinTemperature() && d.EutecticT >= d.Liquidus.MaxTemperature() {
		return fmt.Errorf("eutectic temperature %v K must be below both end-member melting points", d.EutecticT)
	}
	if err := ValidateComposition(d.AlphaMax); err != nil {
		return fmt.Errorf("alpha max solubility: %w", err)
	}
	if err := ValidateComposition(d.BetaMin); err != nil {
		return fmt.Errorf("beta min solubility: %w", err)
	}
	return nil
}

// CheckCrossing 报告两条折线在给定成分处是否上下一致（供测试盯温度轴）。
// left 返回下方折线温度，right 返回上方折线温度。
func CheckCrossing(lo, hi Polyline, c float64) (float64, float64, error) {
	tLo, err := lo.TemperatureAt(c)
	if err != nil {
		return 0, 0, err
	}
	tHi, err := hi.TemperatureAt(c)
	if err != nil {
		return 0, 0, err
	}
	if tLo > tHi {
		return tLo, tHi, fmt.Errorf("polyline order reversed at c=%v: lo=%v K hi=%v K", c, tLo, tHi)
	}
	return tLo, tHi, nil
}
