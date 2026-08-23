package web

import (
	"encoding/json"

	"lever-phase/internal/lever"
)

// pointRequest 是 POST /api/point 的请求体：合金成分 c（质量分数）
// 与温度 t（开尔文）。
type pointRequest struct {
	C float64 `json:"c"`
	T float64 `json:"t"`
}

// scanRequest 是 POST /api/scan 的请求体：固定成分 c，在
// [tmin, tmax] 上扫 n 个温度点。
type scanRequest struct {
	C    float64 `json:"c"`
	TMin float64 `json:"tmin"`
	TMax float64 `json:"tmax"`
	N    int     `json:"n"`
}

// pointResponse 是 /api/point 的响应体，字段与 lever.PointResult
// 保持一致，并额外给出液相分数与总分数便于页面核对守恒。
type pointResponse struct {
	Region          string                `json:"region"`
	RegionLabel     string                `json:"region_label"`
	Phases          []lever.PhaseFraction `json:"phases"`
	EutecticFraction *float64             `json:"eutectic_fraction,omitempty"`
	Note            string                `json:"note"`
	LiquidFraction  float64               `json:"liquid_fraction"`
	TotalFraction   float64               `json:"total_fraction"`
}

// toPointResponse 把领域结果映射为 API 响应。
func toPointResponse(res lever.PointResult) pointResponse {
	return pointResponse{
		Region:          string(res.Region),
		RegionLabel:     lever.ExplainRegion(res.Region),
		Phases:          res.Phases,
		EutecticFraction: res.EutecticFraction,
		Note:            res.Note,
		LiquidFraction:  res.LiquidFraction(),
		TotalFraction:   res.TotalFraction(),
	}
}

// scanResponse 是 /api/scan 的响应体。
type scanResponse struct {
	Composition float64           `json:"composition"`
	TEutecticK  float64           `json:"teutectic_k"`
	Points      []lever.ScanPoint `json:"points"`
}

// examplesResponse 是 GET /api/examples 的响应体：键为算例名，
// 值为算例原始 JSON（页面可一键加载）。
type examplesResponse struct {
	Alloy60 json.RawMessage `json:"alloy-60"`
}
