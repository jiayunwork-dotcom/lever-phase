package phasediag

import "fmt"

// PhaseDiagram 是二元合金简化相图的完整描述：折线液相线、α/β
// 固相线以及共晶参数。所有成分一律为质量分数（0～1）。
type PhaseDiagram struct {
	// Liquidus 是 V 形液相线折线，结点为
	// (0, TAlphaMelting) → (cE, TE) → (1, TBetaMelting)。
	Liquidus Polyline
	// SolidusAlpha 是 α 固溶体边界折线，
	// (0, TAlphaMelting) → (cAlphaMax, TE)。
	SolidusAlpha Polyline
	// SolidusBeta 是 β 固溶体边界折线，
	// (cBetaMin, TE) → (1, TBetaMelting)。
	SolidusBeta Polyline

	EutecticT float64
	EutecticC float64
	AlphaMax  float64
	BetaMin   float64
}

// DefaultPbSn 返回经典 Pb–Sn 型简化相图：固溶 + 共晶、端元熔点与
// 共晶参数钉在 constants.go 的常数上。相图在构造时不重新校验，
// 调用方可用 Validate 确认拓扑一致性。
func DefaultPbSn() PhaseDiagram {
	return PhaseDiagram{
		Liquidus: Polyline{nodes: []Node{
			{Composition: 0, Temperature: TAlphaMelting},
			{Composition: CEutectic, Temperature: TEutectic},
			{Composition: 1, Temperature: TBetaMelting},
		}},
		SolidusAlpha: Polyline{nodes: []Node{
			{Composition: 0, Temperature: TAlphaMelting},
			{Composition: CAlphaMax, Temperature: TEutectic},
		}},
		SolidusBeta: Polyline{nodes: []Node{
			{Composition: CBetaMin, Temperature: TEutectic},
			{Composition: 1, Temperature: TBetaMelting},
		}},
		EutecticT: TEutectic,
		EutecticC: CEutectic,
		AlphaMax:  CAlphaMax,
		BetaMin:   CBetaMin,
	}
}

// Validate 校验整张相图：三条折线各自结点有序，且共晶参数与
// 端点熔点在拓扑上一致。任何不一致都返回错误。
func (d PhaseDiagram) Validate() error {
	if err := d.Liquidus.Validate(); err != nil {
		return fmt.Errorf("liquidus: %w", err)
	}
	if err := d.SolidusAlpha.Validate(); err != nil {
		return fmt.Errorf("alpha solidus: %w", err)
	}
	if err := d.SolidusBeta.Validate(); err != nil {
		return fmt.Errorf("beta solidus: %w", err)
	}
	if err := ValidateDiagramTopology(d); err != nil {
		return err
	}
	// 固相线端点与液相线端点在共晶温度处对齐。
	alphaLast := d.SolidusAlpha.NodeAt(d.SolidusAlpha.NumNodes() - 1)
	if !Within(alphaLast.Composition, d.AlphaMax, 1e-9) {
		return fmt.Errorf("alpha solidus last node composition %v != alpha max %v", alphaLast.Composition, d.AlphaMax)
	}
	betaFirst := d.SolidusBeta.NodeAt(0)
	if !Within(betaFirst.Composition, d.BetaMin, 1e-9) {
		return fmt.Errorf("beta solidus first node composition %v != beta min %v", betaFirst.Composition, d.BetaMin)
	}
	return nil
}

// LiquidusTemperature 由合金成分求液相线温度。成分落在左段用
// 左段插值，落在右段用右段插值，恰好为共晶成分时返回共晶温度。
func (d PhaseDiagram) LiquidusTemperature(c float64) (float64, error) {
	if err := ValidateComposition(c); err != nil {
		return 0, err
	}
	return d.Liquidus.TemperatureAt(c)
}

// LiquidusSegmentIndex 判断成分落在液相线的哪一段（0=左段, 1=右段）。
func (d PhaseDiagram) LiquidusSegmentIndex(c float64) int {
	if c >= d.EutecticC {
		return 1
	}
	return 0
}
