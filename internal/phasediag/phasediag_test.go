package phasediag

import (
	"math"
	"testing"
)

func TestDefaultPbSnConsistent(t *testing.T) {
	d := DefaultPbSn()
	if err := d.Validate(); err != nil {
		t.Fatalf("default Pb-Sn diagram must be valid, got error: %v", err)
	}
	// liquidus must cross both end members at their melting points.
	if !Within(mustTemp(d, 0), TAlphaMelting, 1e-9) {
		t.Errorf("liquidus at c=0 = %.6f, want TAlphaMelting %.6f", mustTemp(d, 0), TAlphaMelting)
	}
	if !Within(mustTemp(d, 1), TBetaMelting, 1e-9) {
		t.Errorf("liquidus at c=1 = %.6f, want TBetaMelting %.6f", mustTemp(d, 1), TBetaMelting)
	}
}

func mustTemp(d PhaseDiagram, c float64) float64 {
	t, err := d.LiquidusTemperature(c)
	if err != nil {
		return math.NaN()
	}
	return t
}

func TestLiquidusTemperatureAxis(t *testing.T) {
	d := DefaultPbSn()
	cases := []struct {
		c      float64
		wantT  float64
	}{
		{0, TAlphaMelting},
		{1, TBetaMelting},
		{CEutectic, TEutectic},
		{CEutectic / 2, TAlphaMelting + (TEutectic-TAlphaMelting)*0.5},
	}
	for _, tc := range cases {
		got, err := d.LiquidusTemperature(tc.c)
		if err != nil {
			t.Fatalf("liquidus temperature at c=%v: %v", tc.c, err)
		}
		if math.Abs(got-tc.wantT) > 1e-6 {
			t.Errorf("liquidus at c=%v = %.6f K, want %.6f K", tc.c, got, tc.wantT)
		}
	}
	// Temperature axis must round-trip: T -> composition -> T.
	temp := 500.0
	comp, err := d.LiquidLeftAt(temp)
	if err != nil {
		t.Fatalf("liquidus-left at T=%v: %v", temp, err)
	}
	back, err := d.LiquidusTemperature(comp)
	if err != nil {
		t.Fatalf("liquidus temperature at c=%v: %v", comp, err)
	}
	if math.Abs(back-temp) > 1e-6 {
		t.Errorf("round trip T=%v -> c=%v -> T=%v, want back at T=%v", temp, comp, back, temp)
	}
}

func TestSolidusTemperatureAxis(t *testing.T) {
	d := DefaultPbSn()
	// alpha solidus endpoints: pure A at TA, alpha max at TE.
	ca, err := d.AlphaBoundaryAt(TEutectic)
	if err != nil {
		t.Fatalf("alpha boundary at TE: %v", err)
	}
	if math.Abs(ca-CAlphaMax) > 1e-9 {
		t.Errorf("alpha boundary at TE = %.6f, want CAlphaMax %.6f", ca, CAlphaMax)
	}
	ca0, err := d.AlphaBoundaryAt(TAlphaMelting)
	if err != nil {
		t.Fatalf("alpha boundary at TA: %v", err)
	}
	if math.Abs(ca0) > 1e-9 {
		t.Errorf("alpha boundary at TA = %.6f, want 0", ca0)
	}
	// beta solidus endpoints: beta min at TE, pure B at TB.
	cb, err := d.BetaBoundaryAt(TEutectic)
	if err != nil {
		t.Fatalf("beta boundary at TE: %v", err)
	}
	if math.Abs(cb-CBetaMin) > 1e-9 {
		t.Errorf("beta boundary at TE = %.6f, want CBetaMin %.6f", cb, CBetaMin)
	}
	cb1, err := d.BetaBoundaryAt(TBetaMelting)
	if err != nil {
		t.Fatalf("beta boundary at TB: %v", err)
	}
	if math.Abs(cb1-1) > 1e-9 {
		t.Errorf("beta boundary at TB = %.6f, want 1", cb1)
	}
	// A mid-range value on the alpha solidus must satisfy the straight line.
	mid := 500.0
	got, err := d.AlphaBoundaryAt(mid)
	if err != nil {
		t.Fatalf("alpha boundary at T=%v: %v", mid, err)
	}
	want := CAlphaMax * (TAlphaMelting - mid) / (TAlphaMelting - TEutectic)
	if math.Abs(got-want) > 1e-6 {
		t.Errorf("alpha boundary at T=%v = %.6f, want straight-line %.6f", mid, got, want)
	}
}

func TestValidateDiagramRejectsUnsortedNodes(t *testing.T) {
	// Decreasing node composition is node disorder and must be rejected.
	bad := []Node{{0.2, 500}, {0.1, 460}, {0.3, 480}}
	if err := ValidatePolylineNodes(bad); err == nil {
		t.Error("polyline with decreasing composition must be rejected")
	}
	// Multiple temperature direction turns (valley-peak-valley) are disorder.
	bad = []Node{{0, 600}, {0.3, 500}, {0.5, 540}, {0.7, 470}, {1, 520}}
	if err := ValidatePolylineNodes(bad); err == nil {
		t.Error("polyline with multiple temperature turns must be rejected")
	}
	// A single-valley liquidus is legitimate.
	ok := []Node{{0, 600.61}, {0.5, 456}, {1, 505.08}}
	if err := ValidatePolylineNodes(ok); err != nil {
		t.Errorf("single-valley liquidus must pass validation, got: %v", err)
	}
	// Flat temperature segment is disorder too.
	flat := []Node{{0, 600}, {0.4, 456}, {0.5, 456}, {1, 505}}
	if err := ValidatePolylineNodes(flat); err == nil {
		t.Error("polyline with flat temperature segment must be rejected")
	}
}

func TestPolylineRootBisection(t *testing.T) {
	p, err := NewPolyline([]Node{{0, TAlphaMelting}, {CEutectic, TEutectic}, {1, TBetaMelting}})
	if err != nil {
		t.Fatalf("build liquidus polyline: %v", err)
	}
	for _, temp := range []float64{500.0, 470.0} {
		roots, err := p.CompositionAtAll(temp)
		if err != nil {
			t.Fatalf("composition at T=%v: %v", temp, err)
		}
		// Between TE and both melting points the V-shaped liquidus has two roots.
		if len(roots) != 2 {
			t.Errorf("T=%v: %d roots, want 2", temp, len(roots))
		}
		for _, c := range roots {
			back, err := p.TemperatureAt(c)
			if err != nil {
				t.Fatalf("temperature at c=%v: %v", c, err)
			}
			if math.Abs(back-temp) > 1e-6 {
				t.Errorf("round trip c=%v: T=%v, want %v", c, back, temp)
			}
		}
	}
}
