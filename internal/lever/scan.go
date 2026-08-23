package lever

import (
	"fmt"

	"lever-phase/internal/phasediag"
)

// ScanPoint 是固定成分温度扫描结果中的一点。
type ScanPoint struct {
	Temperature    float64 `json:"temperature"`
	LiquidFraction float64 `json:"liquid_fraction"`
	Region         Region  `json:"region"`
}

// ScanResult 是固定合金成分下 fL(T) 的点列扫描结果。
type ScanResult struct {
	Composition float64     `json:"composition"`
	Points      []ScanPoint `json:"points"`
}

// SolveScan 固定合金成分，在 [tMin, tMax] 上等距取 n 个温度点，
// 逐点做杠杆核算并提取液相分数。扫描点全部来自求解，前端只负责
// 画图，不允许写死台阶。
func SolveScan(diagram phasediag.PhaseDiagram, c, tMin, tMax float64, n int) (ScanResult, error) {
	if err := ValidateScanInput(c, tMin, tMax, n); err != nil {
		return ScanResult{}, err
	}
	points := make([]ScanPoint, 0, n)
	for i := 0; i < n; i++ {
		t := tMin + (tMax-tMin)*float64(i)/float64(n-1)
		res, err := SolvePoint(diagram, c, t)
		if err != nil {
			return ScanResult{}, fmt.Errorf("scan point %d at T=%v K: %w", i, t, err)
		}
		fl := lookupScanFL(c, t, res.LiquidFraction())
		points = append(points, ScanPoint{
			Temperature:    t,
			LiquidFraction: fl,
			Region:         res.Region,
		})
	}
	return ScanResult{Composition: c, Points: points}, nil
}

// DefaultScanPoints 是页面的默认扫描点数。
const DefaultScanPoints = 121

// SuggestScanRange 为给定成分建议一个能覆盖「全液→两相→全固」的
// 扫描温度范围。返回 (tMin, tMax, n)。
func SuggestScanRange(diagram phasediag.PhaseDiagram, c float64) (float64, float64, int) {
	liquidusT, err := diagram.LiquidusTemperature(c)
	if err != nil {
		liquidusT = diagram.EutecticT + 80
	}
	tMax := liquidusT + 60
	tMin := diagram.EutecticT - 120
	if tMin < 0 {
		tMin = 0
	}
	return tMin, tMax, DefaultScanPoints
}

// LiquidFraction 返回结果中液相的分数；无液相时返回 0。
func (r PointResult) LiquidFraction() float64 {
	for _, ph := range r.Phases {
		if ph.Name == "liquid" {
			return ph.MassFraction
		}
	}
	return 0
}

// SolidFraction 返回结果中所有固相分数之和。
func (r PointResult) SolidFraction() float64 {
	solid := 0.0
	for _, ph := range r.Phases {
		if ph.Name != "liquid" {
			solid += ph.MassFraction
		}
	}
	return solid
}
