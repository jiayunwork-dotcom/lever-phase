package phasediag

// EutecticParams 汇总共晶相关的钉定参数。
type EutecticParams struct {
	Temperature float64
	Composition float64
	AlphaMax    float64
	BetaMin     float64
}

// Eutectic 返回共晶参数快照。
func (d PhaseDiagram) Eutectic() EutecticParams {
	return EutecticParams{
		Temperature: d.EutecticT,
		Composition: d.EutecticC,
		AlphaMax:    d.AlphaMax,
		BetaMin:     d.BetaMin,
	}
}

// InEutecticRange 报告成分是否落在共晶凝固区间内
// （[cAlphaMax, cBetaMin]，温度低于共晶温度时出现 α+β 两相固）。
func (d PhaseDiagram) InEutecticRange(c float64) bool {
	return c >= d.AlphaMax-CompTol && c <= d.BetaMin+CompTol
}

// BelowEutectic 报告温度是否不高于共晶温度。
func (d PhaseDiagram) BelowEutectic(t float64) bool {
	return t <= d.EutecticT+TempTol
}

// AboveEutectic 报告温度是否高于共晶温度。
func (d PhaseDiagram) AboveEutectic(t float64) bool {
	return t > d.EutecticT+TempTol
}

// EutecticFractionWidth 返回共晶凝固区间的成分宽度（cBetaMin − cAlphaMax）。
func (d PhaseDiagram) EutecticFractionWidth() float64 {
	return d.BetaMin - d.AlphaMax
}
