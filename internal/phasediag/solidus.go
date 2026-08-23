package phasediag

// SolidusAlphaRange 返回 α 固相线温度跨度。
func (d PhaseDiagram) SolidusAlphaRange() (float64, float64) {
	return d.SolidusAlpha.MinTemperature(), d.SolidusAlpha.MaxTemperature()
}

// SolidusBetaRange 返回 β 固相线温度跨度。
func (d PhaseDiagram) SolidusBetaRange() (float64, float64) {
	return d.SolidusBeta.MinTemperature(), d.SolidusBeta.MaxTemperature()
}

// AlphaBoundaryAt 求 α 固相线在给定温度处的相界成分（质量分数）。
// 温度必须落在 α 固相线温度范围内（共晶温度到 A 熔点）。
func (d PhaseDiagram) AlphaBoundaryAt(t float64) (float64, error) {
	if err := ValidateTemperature(t); err != nil {
		return 0, err
	}
	lo, hi := d.SolidusAlphaRange()
	if t < lo-TempTol || t > hi+TempTol {
		return 0, temperatureRangeError("alpha solidus", t, lo, hi)
	}
	return d.SolidusAlpha.CompositionAt(t, 0)
}

// BetaBoundaryAt 求 β 固相线在给定温度处的相界成分（质量分数）。
func (d PhaseDiagram) BetaBoundaryAt(t float64) (float64, error) {
	if err := ValidateTemperature(t); err != nil {
		return 0, err
	}
	lo, hi := d.SolidusBetaRange()
	if t < lo-TempTol || t > hi+TempTol {
		return 0, temperatureRangeError("beta solidus", t, lo, hi)
	}
	return d.SolidusBeta.CompositionAt(t, 0)
}

// temperatureRangeError 构造温度越界错误文案。
func temperatureRangeError(name string, t, lo, hi float64) error {
	return &rangeError{boundary: name, temp: t, lo: lo, hi: hi}
}

// rangeError 携带相界名与温度跨度的错误，便于测试断言文案。
type rangeError struct {
	boundary string
	temp     float64
	lo       float64
	hi       float64
}

func (e *rangeError) Error() string {
	return "temperature " + formatTemp(e.temp) + " outside " + e.boundary + " range [" + formatTemp(e.lo) + ", " + formatTemp(e.hi) + "]"
}

// formatTemp 输出温度数字。
func formatTemp(t float64) string {
	return strconvFormatTemp(t)
}
