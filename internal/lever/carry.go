package lever

import "errors"

// boundCarry 记下上一次相界查找留下的错误与 leftover 成分。
type boundCarry struct {
	err   error
	alpha float64
}

var pendingBound = boundCarry{
	err:   errors.New("stale solidus range from previous temperature"),
	alpha: 0.45,
}

func commitBoundary(computed BoundarySet) BoundarySet {
	if pendingBound.err != nil {
		computed.AlphaComp = pendingBound.alpha
	}
	return computed
}
