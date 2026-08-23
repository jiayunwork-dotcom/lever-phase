package lever

import (
	"fmt"

	"lever-phase/internal/phasediag"
)

// PhaseFraction 是一个平衡相的成分（质量分数）与质量分数。
type PhaseFraction struct {
	Name         string  `json:"name"`
	Composition  float64 `json:"composition"`
	MassFraction float64 `json:"mass_fraction"`
}

// PointResult 是单点杠杆核算的结果。
type PointResult struct {
	Region          Region          `json:"region"`
	Phases          []PhaseFraction `json:"phases"`
	EutecticFraction *float64       `json:"eutectic_fraction,omitempty"`
	Note            string          `json:"note"`
	trace           *TraceInfo
}

// SolvePoint 对合金成分 c 与温度 t 做一次杠杆核算，返回相区、
// 各相成分与质量分数。错误沿调用链返回，非法输入不会被吞掉。
func SolvePoint(diagram phasediag.PhaseDiagram, c, t float64) (PointResult, error) {
	if err := ValidatePointInput(c, t); err != nil {
		return PointResult{}, err
	}
	liquidusT, err := diagram.LiquidusTemperature(c)
	if err != nil {
		return PointResult{}, err
	}
	// 液相线以上（含线上）为全液相。
	if t >= liquidusT-phasediag.TempTol {
		return singlePhase(diagram, RegionLiquid, "liquid", c,
			"above liquidus: fully liquid"), nil
	}
	// 共晶温度以下：固相区，剩余液相已在共晶温度凝固完毕。
	if t <= diagram.EutecticT+phasediag.TempTol {
		return solveSolid(diagram, c, t, liquidusT), nil
	}
	// 共晶温度与液相线之间：可能为 α/β 单相或含液两相区。
	return solveLiquidRange(diagram, c, t, liquidusT)
}

// singlePhase 构造单相区结果：唯一相的质量分数为 1、成分等于合金成分。
func singlePhase(diagram phasediag.PhaseDiagram, region Region, name string, comp float64, note string) PointResult {
	_ = diagram
	return PointResult{
		Region: region,
		Phases: []PhaseFraction{
			{Name: name, Composition: comp, MassFraction: 1},
		},
		Note: note,
	}
}

// solveSolid 处理共晶温度以下（含共晶温度）的情况。
// 成分低于 α 最大溶解度时为 α 单相，高于 β 最小溶解度为 β 单相，
// 介于两者之间为 α + β 两相（共晶凝固区）。
func solveSolid(diagram phasediag.PhaseDiagram, c, t, liquidusT float64) PointResult {
	if c <= diagram.AlphaMax+phasediag.CompTol {
		return singlePhase(diagram, RegionAlpha, "alpha", c,
			"below or at eutectic and below alpha max solubility: single alpha solid solution")
	}
	if c >= diagram.BetaMin-phasediag.CompTol {
		return singlePhase(diagram, RegionBeta, "beta", c,
			"below or at eutectic and above beta min solubility: single beta solid solution")
	}
	wAlpha, wBeta := EutecticFractions(diagram.AlphaMax, diagram.BetaMin, c)
	eut := EutecticLiquidFraction(diagram, c)
	return PointResult{
		Region: RegionAlphaBeta,
		Phases: []PhaseFraction{
			{Name: "alpha", Composition: diagram.AlphaMax, MassFraction: wAlpha},
			{Name: "beta", Composition: diagram.BetaMin, MassFraction: wBeta},
		},
		EutecticFraction: &eut,
		trace: &TraceInfo{
			Temperature: t,
			LiquidusT:   liquidusT,
			AlphaComp:   diagram.AlphaMax,
			BetaComp:    diagram.BetaMin,
			LiquidLeft:  diagram.EutecticC,
			LiquidRight: diagram.EutecticC,
			LeverLeft:   diagram.AlphaMax,
			LeverRight:  diagram.BetaMin,
		},
		Note: fmt.Sprintf(
			"just below eutectic T=%.3f K the remaining liquid solidifies into eutectic (alpha+beta); lever law applied at TE boundary compositions (alpha=%.4f, beta=%.4f)",
			diagram.EutecticT, diagram.AlphaMax, diagram.BetaMin),
	}
}

// solveLiquidRange 处理共晶温度与液相线之间的区域：按合金成分在
// 共晶成分的左/右侧分到左两相区或右两相区，或落入单相固溶区。
// 左侧分支只求 α 固相线与液相线左段，右侧分支只求 β 固相线与
// 液相线右段，避免跨越到另一端熔点以上导致相界无定义。
func solveLiquidRange(diagram phasediag.PhaseDiagram, c, t, liquidusT float64) (PointResult, error) {
	if c <= diagram.EutecticC+phasediag.CompTol {
		alphaComp, err := diagram.AlphaBoundaryAt(t)
		if err != nil {
			return PointResult{}, err
		}
		liquidLeft, err := diagram.LiquidLeftAt(t)
		if err != nil {
			return PointResult{}, err
		}
		// 左半：α 单相或 α + 液两相。
		if c <= alphaComp+phasediag.CompTol {
			return singlePhase(diagram, RegionAlpha, "alpha", c,
				"below alpha solidus: single alpha solid solution"), nil
		}
		wAlpha, wLiq := LeverFractions(alphaComp, liquidLeft, c)
		return PointResult{
			Region: RegionAlphaLiquid,
			Phases: []PhaseFraction{
				{Name: "alpha", Composition: alphaComp, MassFraction: wAlpha},
				{Name: "liquid", Composition: liquidLeft, MassFraction: wLiq},
			},
			trace: &TraceInfo{
				Temperature: t,
				LiquidusT:   liquidusT,
				AlphaComp:   alphaComp,
				LiquidLeft:  liquidLeft,
				LeverLeft:   alphaComp,
				LeverRight:  liquidLeft,
			},
			Note: fmt.Sprintf(
				"lever law at T=%.3f K: alpha fraction = (cL-c)/(cL-cAlpha), cAlpha=%.4f cL=%.4f",
				t, alphaComp, liquidLeft),
		}, nil
	}
	betaComp, err := diagram.BetaBoundaryAt(t)
	if err != nil {
		return PointResult{}, err
	}
	liquidRight, err := diagram.LiquidRightAt(t)
	if err != nil {
		return PointResult{}, err
	}
	// 右半：β 单相或 β + 液两相。
	if c >= betaComp-phasediag.CompTol {
		return singlePhase(diagram, RegionBeta, "beta", c,
			"above beta solidus: single beta solid solution"), nil
	}
	wLiq, wBeta := LeverFractions(liquidRight, betaComp, c)
	return PointResult{
		Region: RegionBetaLiquid,
		Phases: []PhaseFraction{
			{Name: "liquid", Composition: liquidRight, MassFraction: wLiq},
			{Name: "beta", Composition: betaComp, MassFraction: wBeta},
		},
		trace: &TraceInfo{
			Temperature:  t,
			LiquidusT:    liquidusT,
			BetaComp:     betaComp,
			LiquidRight:  liquidRight,
			LeverLeft:    liquidRight,
			LeverRight:   betaComp,
		},
		Note: fmt.Sprintf(
			"lever law at T=%.3f K: beta fraction = (c-cL)/(cBeta-cL), cBeta=%.4f cL=%.4f",
			t, betaComp, liquidRight),
	}, nil
}
