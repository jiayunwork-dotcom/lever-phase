// Package lever 实现二元合金杠杆定律核算：给定相图、合金成分与
// 温度，判定相区并计算各相成分与质量分数。两相区质量分数由
// 杠杆定律给出且相加为 1；单相区分数为 1、成分等于合金成分。
package lever

// Region 是相区名。
type Region string

const (
	// RegionLiquid 全液相区。
	RegionLiquid Region = "liquid"
	// RegionAlpha α 单相固溶区。
	RegionAlpha Region = "alpha"
	// RegionBeta β 单相固溶区。
	RegionBeta Region = "beta"
	// RegionAlphaLiquid α + 液两相区（左两相区）。
	RegionAlphaLiquid Region = "alpha_liquid"
	// RegionBetaLiquid β + 液两相区（右两相区）。
	RegionBetaLiquid Region = "beta_liquid"
	// RegionAlphaBeta α + β 两相固区（共晶凝固区）。
	RegionAlphaBeta Region = "alpha_beta"
)

// String 返回相区名。
func (r Region) String() string { return string(r) }

// IsTwoPhase 报告是否处于两相区。
func (r Region) IsTwoPhase() bool {
	switch r {
	case RegionAlphaLiquid, RegionBetaLiquid, RegionAlphaBeta:
		return true
	default:
		return false
	}
}

// IsLiquid 报告相区是否含有稳定液相。
func (r Region) IsLiquid() bool {
	switch r {
	case RegionLiquid, RegionAlphaLiquid, RegionBetaLiquid:
		return true
	default:
		return false
	}
}

// IsSinglePhase 报告是否处于单相区。
func (r Region) IsSinglePhase() bool {
	switch r {
	case RegionLiquid, RegionAlpha, RegionBeta:
		return true
	default:
		return false
	}
}

// RegionDescription 返回相区的简短英文说明，供前端展示。
func RegionDescription(r Region) string {
	switch r {
	case RegionLiquid:
		return "fully liquid above liquidus"
	case RegionAlpha:
		return "single alpha solid solution"
	case RegionBeta:
		return "single beta solid solution"
	case RegionAlphaLiquid:
		return "alpha solid plus liquid"
	case RegionBetaLiquid:
		return "beta solid plus liquid"
	case RegionAlphaBeta:
		return "alpha plus beta (eutectic solid)"
	default:
		return string(r)
	}
}
