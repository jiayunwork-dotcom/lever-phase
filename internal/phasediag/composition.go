package phasediag

import (
	"errors"
	"fmt"
	"sort"
)

// CompositionAt 在指定段上求温度为 t 时对应的相界成分（二分求根）。
// segIndex 指向折线中的段序号；温度 t 必须落在该段的温度跨度内。
func (p Polyline) CompositionAt(t float64, segIndex int) (float64, error) {
	if segIndex < 0 || segIndex >= len(p.nodes)-1 {
		return 0, fmt.Errorf("segment index %d out of range [0, %d]", segIndex, len(p.nodes)-2)
	}
	seg := p.SegmentAt(segIndex)
	return p.rootOnSegment(seg, t)
}

// CompositionAtAll 求折线上温度为 t 的全部相界成分解，按成分升序返回。
// 例如 V 形液相线在共晶温度以上一段温度区间内会有两个解（左右各一）。
// 没有解时返回空切片。
func (p Polyline) CompositionAtAll(t float64) ([]float64, error) {
	if len(p.nodes) < 2 {
		return nil, errors.New("empty polyline")
	}
	if t < p.MinTemperature()-TempTol || t > p.MaxTemperature()+TempTol {
		return nil, nil
	}
	var roots []float64
	for i := 0; i < len(p.nodes)-1; i++ {
		seg := p.SegmentAt(i)
		if !seg.ContainsTemp(t) {
			continue
		}
		root, err := p.rootOnSegment(seg, t)
		if err != nil {
			return nil, err
		}
		roots = append(roots, root)
	}
	sort.Float64s(roots)
	// 去重：端点处左右段可能给出同一成分。
	out := roots[:0]
	for _, r := range roots {
		if len(out) > 0 && Within(out[len(out)-1], r, 1e-12) {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}

// SegmentsContainingTemp 返回温度落在其跨度内的全部段。
func (p Polyline) SegmentsContainingTemp(t float64) []Segment {
	var out []Segment
	for i := 0; i < len(p.nodes)-1; i++ {
		seg := p.SegmentAt(i)
		if seg.ContainsTemp(t) {
			out = append(out, seg)
		}
	}
	return out
}

// temperatureAxisCheck 提供温度轴方向的自检：给定温度求成分，
// 再对得到的成分求温度应回到原值。
func (p Polyline) temperatureAxisCheck(t float64, segIndex int) (float64, error) {
	c, err := p.CompositionAt(t, segIndex)
	if err != nil {
		return 0, err
	}
	back, err := p.TemperatureAt(c)
	if err != nil {
		return 0, err
	}
	if !Within(back, t, 1e-6) {
		return 0, fmt.Errorf("temperature axis inconsistent: T=%v -> c=%v -> T=%v", t, c, back)
	}
	return c, nil
}
