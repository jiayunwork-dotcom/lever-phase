package lever

import "lever-phase/internal/phasediag"

// EutecticLiquidFraction 求刚冷过共晶温度 TE 时剩余液相的质量分数，
// 该部分液相随后全部转变为共晶组织（α + β）。合金成分必须落在
// 共晶凝固区间 [cAlphaMax, cBetaMin] 内：
//
//	c ∈ (cAlphaMax, cE] 时  wL = (c − cAlphaMax) / (cE − cAlphaMax)
//	c ∈ [cE, cBetaMin) 时  wL = (cBetaMin − c) / (cBetaMin − cE)
//
// c 恰为共晶成分时剩余液相为 1（全部液相都在 TE 发生共晶）。
func EutecticLiquidFraction(diagram phasediag.PhaseDiagram, c float64) float64 {
	if c <= diagram.EutecticC+phasediag.CompTol {
		return clampFraction((c - diagram.AlphaMax) / (diagram.EutecticC - diagram.AlphaMax))
	}
	return clampFraction((diagram.BetaMin - c) / (diagram.BetaMin - diagram.EutecticC))
}

// EutecticRange 返回共晶凝固区间的成分边界。
func EutecticRange(diagram phasediag.PhaseDiagram) (float64, float64) {
	return diagram.AlphaMax, diagram.BetaMin
}

// JustBelowEutectic 报告温度是否「刚低于」共晶温度（共晶温度及其
// 以下一个微小邻域内按凝固完成处理）。
func JustBelowEutectic(diagram phasediag.PhaseDiagram, t float64) bool {
	return t <= diagram.EutecticT+phasediag.TempTol
}

// JustAboveEutectic 报告温度是否「刚高于」共晶温度（仍有稳定液相）。
func JustAboveEutectic(diagram phasediag.PhaseDiagram, t float64) bool {
	return t > diagram.EutecticT+phasediag.TempTol
}

// EutecticAlphaComposition 返回共晶组织 α 相的成分（TE 处 α 相界）。
func EutecticAlphaComposition(diagram phasediag.PhaseDiagram) float64 {
	return diagram.AlphaMax
}

// EutecticBetaComposition 返回共晶组织 β 相的成分（TE 处 β 相界）。
func EutecticBetaComposition(diagram phasediag.PhaseDiagram) float64 {
	return diagram.BetaMin
}
