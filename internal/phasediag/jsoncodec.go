package phasediag

import (
	"encoding/json"
	"fmt"
)

// DiagramJSON 是相图的 JSON 表示，允许调用方用自定义结点喂给内核。
// 解析后必须通过结点有序与拓扑一致性校验。
type DiagramJSON struct {
	Liquidus     []Node  `json:"liquidus"`
	SolidusAlpha []Node  `json:"solidus_alpha"`
	SolidusBeta  []Node  `json:"solidus_beta"`
	EutecticT    float64 `json:"eutectic_t"`
	EutecticC    float64 `json:"eutectic_c"`
	AlphaMax     float64 `json:"alpha_max"`
	BetaMin      float64 `json:"beta_min"`
}

// MarshalJSON 把相图序列化为可读 JSON。
func (d PhaseDiagram) MarshalJSON() ([]byte, error) {
	return json.Marshal(DiagramJSON{
		Liquidus:     d.Liquidus.Nodes(),
		SolidusAlpha: d.SolidusAlpha.Nodes(),
		SolidusBeta:  d.SolidusBeta.Nodes(),
		EutecticT:    d.EutecticT,
		EutecticC:    d.EutecticC,
		AlphaMax:     d.AlphaMax,
		BetaMin:      d.BetaMin,
	})
}

// UnmarshalJSON 从 JSON 解析相图，并用与构造一致的方式完成校验。
func (d *PhaseDiagram) UnmarshalJSON(data []byte) error {
	parsed, err := ParseDiagram(data)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// ParseDiagram 解析 JSON 描述的相图并校验：三条折线结点有序、
// 共晶参数与端点熔点拓扑一致。任何乱序或越界都返回错误。
func ParseDiagram(data []byte) (PhaseDiagram, error) {
	var j DiagramJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return PhaseDiagram{}, fmt.Errorf("invalid diagram JSON: %w", err)
	}
	return BuildDiagram(j)
}

// BuildDiagram 由显式结点构造相图并完成全部校验。
func BuildDiagram(j DiagramJSON) (PhaseDiagram, error) {
	liquidus, err := NewPolyline(j.Liquidus)
	if err != nil {
		return PhaseDiagram{}, fmt.Errorf("liquidus: %w", err)
	}
	alpha, err := NewPolyline(j.SolidusAlpha)
	if err != nil {
		return PhaseDiagram{}, fmt.Errorf("alpha solidus: %w", err)
	}
	beta, err := NewPolyline(j.SolidusBeta)
	if err != nil {
		return PhaseDiagram{}, fmt.Errorf("beta solidus: %w", err)
	}
	d := PhaseDiagram{
		Liquidus:     liquidus,
		SolidusAlpha: alpha,
		SolidusBeta:  beta,
		EutecticT:    j.EutecticT,
		EutecticC:    j.EutecticC,
		AlphaMax:     j.AlphaMax,
		BetaMin:      j.BetaMin,
	}
	if err := d.Validate(); err != nil {
		return PhaseDiagram{}, err
	}
	return d, nil
}

// MarshalDefault 返回默认 Pb–Sn 相图的 JSON 文本。
func MarshalDefault() ([]byte, error) {
	return json.Marshal(DefaultPbSn())
}

// ParseDefault 从 JSON 重建默认相图（用于验证编解码往返一致）。
func ParseDefault(data []byte) (PhaseDiagram, error) {
	return ParseDiagram(data)
}
