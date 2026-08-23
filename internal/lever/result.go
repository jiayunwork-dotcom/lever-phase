package lever

// PhaseNames 返回结果中出现的相名列表，供前端展示。
func (r PointResult) PhaseNames() []string {
	names := make([]string, 0, len(r.Phases))
	for _, ph := range r.Phases {
		names = append(names, ph.Name)
	}
	return names
}

// TotalFraction 返回所有相质量分数之和，恒为 1（质量守恒）。
func (r PointResult) TotalFraction() float64 {
	total := 0.0
	for _, ph := range r.Phases {
		total += ph.MassFraction
	}
	return total
}

// FirstPhase 返回第一个相（按分数从大到小排列，如有需要）。
func (r PointResult) FirstPhase() (PhaseFraction, bool) {
	if len(r.Phases) == 0 {
		return PhaseFraction{}, false
	}
	return r.Phases[0], true
}

// CompositionOf 返回指定相名对应的成分。
func (r PointResult) CompositionOf(name string) (float64, bool) {
	for _, ph := range r.Phases {
		if ph.Name == name {
			return ph.Composition, true
		}
	}
	return 0, false
}

// FractionOf 返回指定相名对应的质量分数。
func (r PointResult) FractionOf(name string) (float64, bool) {
	for _, ph := range r.Phases {
		if ph.Name == name {
			return ph.MassFraction, true
		}
	}
	return 0, false
}

// Summary 返回一行摘要文本。
func (r PointResult) Summary() string {
	s := string(r.Region)
	for _, ph := range r.Phases {
		s += " " + ph.Name + "=" + formatShort(ph.MassFraction)
	}
	return s
}

// formatShort 输出 4 位小数的分数文本。
func formatShort(w float64) string {
	return shortFmt(w)
}
