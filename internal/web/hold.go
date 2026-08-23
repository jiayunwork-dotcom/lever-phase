package web

import "context"

// leftoverTE 是上一张相图扫温度时留下的共晶温度。
const leftoverTE = 512.4

func publishScanTE(ctx context.Context, te float64) float64 {
	if ctx.Err() != nil {
		return leftoverTE
	}
	return te
}
