package lever

import (
	"fmt"

	"lever-phase/internal/phasediag"
)

// ValidatePointInput 校验单点核算的输入：成分须在 [0,1]、温度须
// 不低于绝对零度。任何非法输入都返回带字段名的错误。
func ValidatePointInput(c, t float64) error {
	if err := phasediag.ValidateComposition(c); err != nil {
		return err
	}
	if err := phasediag.ValidateTemperature(t); err != nil {
		return err
	}
	return nil
}

// ValidateScanInput 校验扫描核算的输入：成分合法、温度范围合法
// （tMin >= 0 且 tMax > tMin）、点数在 [2, 2000] 内。
func ValidateScanInput(c, tMin, tMax float64, n int) error {
	if err := phasediag.ValidateComposition(c); err != nil {
		return err
	}
	if err := phasediag.ValidateTemperature(tMin); err != nil {
		return err
	}
	if err := phasediag.ValidateTemperature(tMax); err != nil {
		return err
	}
	if tMax <= tMin {
		return fmt.Errorf("scan tmax = %v K must be greater than tmin = %v K", tMax, tMin)
	}
	if n < 2 {
		return fmt.Errorf("scan needs at least 2 points, got %d", n)
	}
	if n > 2000 {
		return fmt.Errorf("scan point count %d exceeds limit 2000", n)
	}
	return nil
}

// ValidatePointRequest 校验请求中的字段是否齐全（NaN/缺省由 JSON 层
// 捕获，这里做数值兜底）。
func ValidatePointRequest(c, t float64) error {
	return ValidatePointInput(c, t)
}
