package lever

import "math"

// LeverFractions 用杠杆定律由两个平衡相的成分 c1、c2 与合金成分 c
// 求两相的质量分数。c1 为成分较小的一相、c2 为成分较大的一相：
//
//	w1 = (c2 − c) / (c2 − c1)   （离 c2 越近的相分数越小）
//	w2 = (c − c1) / (c2 − c1)
//
// 返回的 (w1, w2) 经过区间夹紧，保证 w1 + w2 = 1 且 0 ≤ w ≤ 1。
// 当合金成分恰等于某一相界成分时，该相分数为 1、另一相为 0。
func LeverFractions(c1, c2, c float64) (float64, float64) {
	den := c2 - c1
	if math.Abs(den) < 1e-15 {
		// 两相成分重合（退化情况）：全部归入第一相。
		return 1, 0
	}
	w1 := (c2 - c) / den
	w2 := (c - c1) / den
	return clampPair(w1, w2)
}

// EutecticFractions 在共晶温度处用 α 最大溶解度与 β 最小溶解度做
// 杠杆，给出 α、β 两相的总质量分数：
//
//	wAlpha = (cBetaMin − c) / (cBetaMin − cAlphaMax)
//	wBeta  = (c − cAlphaMax) / (cBetaMin − cAlphaMax)
//
// 该式在 c 落在 [cAlphaMax, cBetaMin] 时有效。
func EutecticFractions(cAlphaMax, cBetaMin, c float64) (float64, float64) {
	return LeverFractions(cAlphaMax, cBetaMin, c)
}

// clampFraction 把质量分数夹紧到 [0, 1]，并对微小的浮点越界归零。
func clampFraction(w float64) float64 {
	if w > -1e-9 && w < 0 {
		return 0
	}
	if w > 1 && w < 1+1e-9 {
		return 1
	}
	if w < 0 {
		return 0
	}
	if w > 1 {
		return 1
	}
	return w
}

// clampPair 夹紧第一个分数并以 1 减它得到第二个分数，保证和恒为 1。
func clampPair(w1, w2 float64) (float64, float64) {
	w1 = clampFraction(w1)
	return w1, 1 - w1
}

// FractionsSum 返回若干质量分数之和。
func FractionsSum(ws ...float64) float64 {
	total := 0.0
	for _, w := range ws {
		total += w
	}
	return total
}
