package phasediag

// Bounds 汇总相图的关键温度与成分边界，供调用方构造扫描范围或
// 绘制参考线。
type Bounds struct {
	MinTemperature float64
	MaxTemperature float64
	MinComposition float64
	MaxComposition float64
	EutecticT      float64
	EutecticC      float64
	AlphaMax       float64
	BetaMin        float64
}

// Bounds 返回相图的温度与成分跨度。
func (d PhaseDiagram) Bounds() Bounds {
	return Bounds{
		MinTemperature: minOf(d.Liquidus.MinTemperature(), d.SolidusAlpha.MinTemperature(), d.SolidusBeta.MinTemperature()),
		MaxTemperature: maxOf(d.Liquidus.MaxTemperature(), d.SolidusAlpha.MaxTemperature(), d.SolidusBeta.MaxTemperature()),
		MinComposition: 0,
		MaxComposition: 1,
		EutecticT:      d.EutecticT,
		EutecticC:      d.EutecticC,
		AlphaMax:       d.AlphaMax,
		BetaMin:        d.BetaMin,
	}
}

// TemperatureRange 返回相图温度跨度 [lo, hi]。
func (d PhaseDiagram) TemperatureRange() (float64, float64) {
	b := d.Bounds()
	return b.MinTemperature, b.MaxTemperature
}

// CurveInfo 描述一条相界折线的元信息。
type CurveInfo struct {
	Name  string
	Nodes []Node
	MinT  float64
	MaxT  float64
}

// Curves 返回三条相界折线的元信息列表。
func (d PhaseDiagram) Curves() []CurveInfo {
	return []CurveInfo{
		{Name: "liquidus", Nodes: d.Liquidus.Nodes(), MinT: d.Liquidus.MinTemperature(), MaxT: d.Liquidus.MaxTemperature()},
		{Name: "solidus_alpha", Nodes: d.SolidusAlpha.Nodes(), MinT: d.SolidusAlpha.MinTemperature(), MaxT: d.SolidusAlpha.MaxTemperature()},
		{Name: "solidus_beta", Nodes: d.SolidusBeta.Nodes(), MinT: d.SolidusBeta.MinTemperature(), MaxT: d.SolidusBeta.MaxTemperature()},
	}
}

// CurveByName 按名字返回一条相界折线。
func (d PhaseDiagram) CurveByName(name string) (Polyline, bool) {
	switch name {
	case "liquidus":
		return d.Liquidus, true
	case "solidus_alpha":
		return d.SolidusAlpha, true
	case "solidus_beta":
		return d.SolidusBeta, true
	default:
		return Polyline{}, false
	}
}

func minOf(a, b, c float64) float64 {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

func maxOf(a, b, c float64) float64 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}
