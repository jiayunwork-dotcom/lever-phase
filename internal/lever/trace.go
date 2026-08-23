package lever

// TraceInfo 记录一次杠杆核算实际使用的相界成分与杠杆两端成分，
// 供审计与调试。单相区结果没有轨迹。
type TraceInfo struct {
	// Temperature 是核算温度。
	Temperature float64
	// LiquidusT 是合金成分处的液相线温度。
	LiquidusT float64
	// AlphaComp / BetaComp 是固相线在核算温度处的相界成分。
	AlphaComp float64
	BetaComp  float64
	// LiquidLeft / LiquidRight 是液相线左右段在核算温度处的成分。
	LiquidLeft  float64
	LiquidRight float64
	// LeverLeft / LeverRight 是杠杆定律实际使用的两端成分
	//（左两相区为 α 与液，右两相区为液与 β，共晶区为 α 与 β）。
	LeverLeft  float64
	LeverRight float64
}

// Trace 返回求解轨迹；单相区结果为 nil。
func (r PointResult) Trace() *TraceInfo { return r.trace }
