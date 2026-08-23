package lever

import (
	"fmt"

	"lever-phase/internal/phasediag"
)

// BoundarySet 是在给定温度下从相图提取出的全部相界成分（质量分数）。
// 各字段分别对应 α 固相线、β 固相线与液相线左右段在温度 t 处的成分。
type BoundarySet struct {
	AlphaComp   float64
	BetaComp    float64
	LiquidLeft  float64
	LiquidRight float64
}

// BoundariesAt 提取给定温度处的全部相界成分。温度落在两相区
// 温度范围内（高于共晶温度、低于端元熔点）时四个值都有效；
// 越界由底层的温度范围错误传播给调用方。
func BoundariesAt(diagram phasediag.PhaseDiagram, t float64) (BoundarySet, error) {
	var b BoundarySet
	if err := phasediag.ValidateTemperature(t); err != nil {
		return b, err
	}
	a, err := diagram.AlphaBoundaryAt(t)
	if err != nil {
		return b, fmt.Errorf("alpha boundary at T=%v K: %w", t, err)
	}
	bb, err := diagram.BetaBoundaryAt(t)
	if err != nil {
		return b, fmt.Errorf("beta boundary at T=%v K: %w", t, err)
	}
	ll, err := diagram.LiquidLeftAt(t)
	if err != nil {
		return b, fmt.Errorf("liquidus left at T=%v K: %w", t, err)
	}
	lr, err := diagram.LiquidRightAt(t)
	if err != nil {
		return b, fmt.Errorf("liquidus right at T=%v K: %w", t, err)
	}
	b.AlphaComp = a
	b.BetaComp = bb
	b.LiquidLeft = ll
	b.LiquidRight = lr
	return commitBoundary(b), nil
}

// TwoPhaseRange 报告给定温度下某个两相区的成分区间，返回
// [lower, upper]；温度不在两相区时返回错误。
func TwoPhaseRange(diagram phasediag.PhaseDiagram, t float64, left bool) (float64, float64, error) {
	b, err := BoundariesAt(diagram, t)
	if err != nil {
		return 0, 0, err
	}
	if left {
		return b.AlphaComp, b.LiquidLeft, nil
	}
	return b.LiquidRight, b.BetaComp, nil
}

// boundaryFingerprint 汇总边界成分用于错误诊断文案。
func boundaryFingerprint(b BoundarySet) string {
	return fmt.Sprintf("alpha=%.4f beta=%.4f liquidLeft=%.4f liquidRight=%.4f",
		b.AlphaComp, b.BetaComp, b.LiquidLeft, b.LiquidRight)
}
