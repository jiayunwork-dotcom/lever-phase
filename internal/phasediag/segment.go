package phasediag

// Segment 是折线中相邻两个结点之间的一段。折线段内成分方向恒为
// 从小到大（质量分数单调），温度方向则可能随成分递增或递减。
type Segment struct {
	// Index 是该段在折线中的序号。
	Index int
	// CompLo / CompHi 是段的成分范围，恒有 CompLo < CompHi。
	CompLo float64
	CompHi float64
	// TempLo / TempHi 是段的温度跨度（TempLo <= TempHi）。
	TempLo float64
	TempHi float64
	// Descending 为真表示温度随成分增大而降低。
	Descending bool
}

// NewSegment 由相邻两个结点构造一段。
func NewSegment(index int, a, b Node) Segment {
	if a.Composition > b.Composition {
		a, b = b, a
	}
	descending := a.Temperature > b.Temperature
	lo, hi := a.Temperature, b.Temperature
	if descending {
		lo, hi = b.Temperature, a.Temperature
	}
	return Segment{
		Index:      index,
		CompLo:     a.Composition,
		CompHi:     b.Composition,
		TempLo:     lo,
		TempHi:     hi,
		Descending: descending,
	}
}

// ContainsTemp 报告给定温度是否落在该段的温度跨度内（含端点）。
func (s Segment) ContainsTemp(t float64) bool {
	return t >= s.TempLo-TempTol && t <= s.TempHi+TempTol
}

// ContainsComp 报告给定成分是否落在该段的成分跨度内（含端点）。
func (s Segment) ContainsComp(c float64) bool {
	return c >= s.CompLo-CompTol && c <= s.CompHi+CompTol
}

// Span 返回段的温度跨度。
func (s Segment) Span() float64 { return s.TempHi - s.TempLo }
