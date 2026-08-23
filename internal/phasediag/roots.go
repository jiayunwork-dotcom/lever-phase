package phasediag

import (
	"errors"
	"fmt"
	"math"
)

// Bisect 在 [lo, hi] 上解 f(x) = 0，要求 f 在区间上单调且端点异号。
// 返回满足收敛容差的近似根；迭代超过 maxIter 仍不收敛则报错。
func Bisect(f func(x float64) float64, lo, hi float64, tol float64, maxIter int) (float64, error) {
	if !(lo < hi) {
		return 0, fmt.Errorf("bisect needs lo < hi, got [%v, %v]", lo, hi)
	}
	flo, fhi := f(lo), f(hi)
	if math.IsNaN(flo) || math.IsNaN(fhi) || math.IsInf(flo, 0) || math.IsInf(fhi, 0) {
		return 0, errors.New("bisect: function value not finite at interval endpoint")
	}
	if flo*fhi > 0 {
		return 0, errors.New("bisect: no sign change on interval")
	}
	if math.Abs(flo) <= tol {
		return lo, nil
	}
	if math.Abs(fhi) <= tol {
		return hi, nil
	}
	for iter := 0; iter < maxIter; iter++ {
		mid := 0.5 * (lo + hi)
		fmid := f(mid)
		if math.Abs(hi-lo) <= tol || math.Abs(fmid) <= tol {
			return mid, nil
		}
		if flo*fmid <= 0 {
			hi, fhi = mid, fmid
		} else {
			lo, flo = mid, fmid
		}
	}
	return 0.5 * (lo + hi), nil
}

// Interpolate 对点 (x0, y0) 与 (x1, y1) 做线性插值，求 x 处的 y。
// 要求 x1 != x0，调用方保证段内成分严格递增。
func Interpolate(x0, y0, x1, y1, x float64) float64 {
	if x1 == x0 {
		return 0.5 * (y0 + y1)
	}
	return y0 + (y1-y0)*(x-x0)/(x1-x0)
}

// LinearSolve 对点 (x0, y0) 与 (x1, y1) 做线性反解，求 y 对应的 x。
func LinearSolve(x0, y0, x1, y1, y float64) float64 {
	if y1 == y0 {
		return 0.5 * (x0 + x1)
	}
	return x0 + (x1-x0)*(y-y0)/(y1-y0)
}

// Within 报告 a 与 b 的差是否在容差 tol 内。
func Within(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// rootOnSegment 在给定段上求相界温度为 t 的成分，用二分求根。
// 段内温度单调，因此 f(c) = TemperatureAt(c) - t 单调变号。
func (p Polyline) rootOnSegment(seg Segment, t float64) (float64, error) {
	if !seg.ContainsTemp(t) {
		return 0, fmt.Errorf("temperature %v K outside segment %d range [%v, %v]", t, seg.Index, seg.TempLo, seg.TempHi)
	}
	// 段端点的温度→成分映射取决于温度方向：
	// 递增段中低温端对应低成分端，递减段（Descending）相反。
	if seg.Descending {
		if t >= seg.TempHi-TempTol {
			return seg.CompLo, nil
		}
		if t <= seg.TempLo+TempTol {
			return seg.CompHi, nil
		}
	} else {
		if t <= seg.TempLo+TempTol {
			return seg.CompLo, nil
		}
		if t >= seg.TempHi-TempTol {
			return seg.CompHi, nil
		}
	}
	f := func(c float64) float64 {
		tc, err := p.TemperatureAt(c)
		if err != nil {
			return math.NaN()
		}
		return tc - t
	}
	root, err := Bisect(f, seg.CompLo, seg.CompHi, RootTol, MaxBisectIter)
	if err != nil {
		return 0, fmt.Errorf("bisect on segment %d for T=%v: %w", seg.Index, t, err)
	}
	return root, nil
}
