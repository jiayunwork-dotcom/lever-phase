// Package phasediag 提供二元合金简化相图的数据模型：折线液相线、
// α/β 固相线与共晶参数，以及沿相界的求根（成分↔温度互求）。
// 成分一律为质量分数（0～1），温度一律为开尔文。
package phasediag

// 默认 Pb–Sn 型相图的钉定常数，取自经典 Pb–Sn 平衡相图的简化折线。
const (
	// TAlphaMelting 是纯 Pb（组元 A）的熔点，单位 K。
	TAlphaMelting = 600.61
	// TBetaMelting 是纯 Sn（组元 B）的熔点，单位 K。
	TBetaMelting = 505.08
	// TEutectic 是共晶温度；低于该温度不再有稳定液相，单位 K。
	TEutectic = 456.0
	// CEutectic 是共晶成分，即共晶液相的 Sn 质量分数。
	CEutectic = 0.619
	// CAlphaMax 是共晶温度处 α 相的最大 Sn 溶解度（质量分数）。
	CAlphaMax = 0.192
	// CBetaMin 是共晶温度处 β 相的最小 Sn 含量（质量分数），
	// 其余为 Pb 溶解度上限。
	CBetaMin = 0.975
)

// 求根与判区的数值容差。
const (
	// RootTol 是相界求根（二分）的收敛容差。
	RootTol = 1e-9
	// TempTol 是相区边界比较用的温度容差。
	TempTol = 1e-7
	// CompTol 是相界成分比较用的容差。
	CompTol = 1e-9
	// MaxBisectIter 是二分求根的最大迭代次数上限。
	MaxBisectIter = 200
)
