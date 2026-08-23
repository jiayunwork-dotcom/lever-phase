package phasediag

import (
	"errors"
	"fmt"
	"math"
)

// Polyline 表示一条成分单调的折线相界，例如液相线或某条固相线。
// 折线结点按成分升序排列；温度在每个段内严格单调，整条折线
// 允许最多一次温度方向翻转（例如液相线在共晶点的 V 形谷）。
type Polyline struct {
	nodes []Node
}

// NewPolyline 校验并构造一条折线。
func NewPolyline(nodes []Node) (Polyline, error) {
	if len(nodes) < 2 {
		return Polyline{}, errors.New("polyline needs at least two nodes")
	}
	p := Polyline{nodes: append([]Node(nil), nodes...)}
	if err := p.Validate(); err != nil {
		return Polyline{}, err
	}
	return p, nil
}

// Nodes 返回结点副本，调用方可以安全修改返回切片而不影响折线。
func (p Polyline) Nodes() []Node {
	return append([]Node(nil), p.nodes...)
}

// NumNodes 返回结点数量。
func (p Polyline) NumNodes() int { return len(p.nodes) }

// NodeAt 返回第 i 个结点。
func (p Polyline) NodeAt(i int) Node { return p.nodes[i] }

// SegmentAt 返回第 i 段（相邻结点之间的折线段）。
func (p Polyline) SegmentAt(i int) Segment {
	return NewSegment(i, p.nodes[i], p.nodes[i+1])
}

// Segments 返回全部段。
func (p Polyline) Segments() []Segment {
	segs := make([]Segment, 0, len(p.nodes)-1)
	for i := 0; i < len(p.nodes)-1; i++ {
		segs = append(segs, p.SegmentAt(i))
	}
	return segs
}

// MinTemperature 返回折线温度跨度下界。
func (p Polyline) MinTemperature() float64 {
	min := p.nodes[0].Temperature
	for _, n := range p.nodes[1:] {
		if n.Temperature < min {
			min = n.Temperature
		}
	}
	return min
}

// MaxTemperature 返回折线温度跨度上界。
func (p Polyline) MaxTemperature() float64 {
	max := p.nodes[0].Temperature
	for _, n := range p.nodes[1:] {
		if n.Temperature > max {
			max = n.Temperature
		}
	}
	return max
}

// TemperatureAt 由给定成分在折线上线性插值出对应的相界温度。
// 成分必须在折线的成分跨度内，否则返回错误。
func (p Polyline) TemperatureAt(comp float64) (float64, error) {
	if len(p.nodes) == 0 {
		return 0, errors.New("empty polyline")
	}
	first, last := p.nodes[0], p.nodes[len(p.nodes)-1]
	if comp < first.Composition-CompTol || comp > last.Composition+CompTol {
		return 0, fmt.Errorf("composition %v outside polyline range [%v, %v]", comp, first.Composition, last.Composition)
	}
	for i := 0; i < len(p.nodes)-1; i++ {
		a, b := p.nodes[i], p.nodes[i+1]
		if comp >= a.Composition-CompTol && comp <= b.Composition+CompTol {
			return Interpolate(a.Composition, a.Temperature, b.Composition, b.Temperature, comp), nil
		}
	}
	// 到达这里说明成分落在两个相邻结点之间但未被捕获（数值角落）。
	return 0, fmt.Errorf("no polyline segment found for composition %v", comp)
}

// temperatureGap 检查折线整体温度是否出现多于一次方向翻转。
// 液相线允许一个 V 形谷；其余折线应为单调。返回翻转次数。
func (p Polyline) temperatureGap() int {
	if len(p.nodes) < 3 {
		return 0
	}
	prevDir := 0
	turns := 0
	for i := 0; i < len(p.nodes)-1; i++ {
		diff := p.nodes[i+1].Temperature - p.nodes[i].Temperature
		var dir int
		switch {
		case diff > 0:
			dir = 1
		case diff < 0:
			dir = -1
		default:
			dir = 0
		}
		if dir == 0 {
			// 水平温度段在相界折线中属于乱序（共晶水平不是折线的一部分）。
			return 2
		}
		if prevDir != 0 && dir != prevDir {
			turns++
		}
		prevDir = dir
	}
	return turns
}

// Validate 校验折线：成分严格递增、温度方向翻转至多一次、温度为有限正值。
func (p Polyline) Validate() error {
	if len(p.nodes) < 2 {
		return errors.New("polyline needs at least two nodes")
	}
	for i, n := range p.nodes {
		if err := ValidateComposition(n.Composition); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
		if err := ValidateTemperature(n.Temperature); err != nil {
			return fmt.Errorf("node %d: %w", i, err)
		}
		if i > 0 && n.Composition <= p.nodes[i-1].Composition {
			return fmt.Errorf("node compositions not strictly increasing at index %d: %s before %s", i, p.nodes[i-1], n)
		}
	}
	if turns := p.temperatureGap(); turns > 1 {
		return fmt.Errorf("node temperatures out of order: %d direction turns in polyline, want at most one valley", turns)
	}
	return nil
}

// sanity 供包内快速判断折线可用性。
func (p Polyline) sanity() bool {
	return len(p.nodes) >= 2 && !math.IsNaN(p.nodes[0].Temperature)
}
