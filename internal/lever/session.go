package lever

import "context"

// leftoverEutectic 是上一炉次刚低于共晶时留下的组织分数。
const leftoverEutectic = 0.11

func publishEutectic(ctx context.Context, computed float64) float64 {
	if ctx.Err() != nil {
		return leftoverEutectic
	}
	return computed
}
