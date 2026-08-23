package lever

import (
	"math"
	"testing"

	"lever-phase/internal/phasediag"
)

func testDiagram(t *testing.T) phasediag.PhaseDiagram {
	t.Helper()
	return phasediag.DefaultPbSn()
}

// TestAlloy60ExampleClosedForm checks the bundled example against the
// lever rule closed form computed directly from the polyline nodes.
func TestAlloy60ExampleClosedForm(t *testing.T) {
	d := testDiagram(t)
	c, T := 0.60, 458.0
	// closed form: straight-line interpolation on alpha solidus and liquidus-left.
	cAlpha := phasediag.CAlphaMax * (phasediag.TAlphaMelting - T) / (phasediag.TAlphaMelting - phasediag.TEutectic)
	cL := phasediag.CEutectic * (phasediag.TAlphaMelting - T) / (phasediag.TAlphaMelting - phasediag.TEutectic)
	wantWLiq := (c - cAlpha) / (cL - cAlpha)
	wantWAlpha := (cL - c) / (cL - cAlpha)

	res, err := SolvePoint(d, c, T)
	if err != nil {
		t.Fatalf("SolvePoint(c=%v, T=%v): %v", c, T, err)
	}
	if res.Region != RegionAlphaLiquid {
		t.Errorf("region = %s, want %s", res.Region, RegionAlphaLiquid)
	}
	if math.Abs(res.LiquidFraction()-wantWLiq) > 1e-6 {
		t.Errorf("liquid fraction = %.6f, closed form want %.6f", res.LiquidFraction(), wantWLiq)
	}
	alphaFrac, ok := res.FractionOf("alpha")
	if !ok {
		t.Fatal("result has no alpha phase")
	}
	if math.Abs(alphaFrac-wantWAlpha) > 1e-6 {
		t.Errorf("alpha fraction = %.6f, closed form want %.6f", alphaFrac, wantWAlpha)
	}
	if math.Abs(res.TotalFraction()-1) > 1e-9 {
		t.Errorf("mass fractions sum to %.9f, want 1", res.TotalFraction())
	}
}

func TestTwoPhaseMassConservation(t *testing.T) {
	d := testDiagram(t)
	points := []struct {
		c float64
		T float64
	}{
		{0.30, 500},
		{0.45, 475},
		{0.60, 458},
		{0.80, 470},
		{0.90, 480},
		{0.40, 400},
		{0.62, 450},
		{0.80, 440},
	}
	for _, p := range points {
		res, err := SolvePoint(d, p.c, p.T)
		if err != nil {
			t.Fatalf("SolvePoint(c=%v, T=%v): %v", p.c, p.T, err)
		}
		if !res.Region.IsTwoPhase() {
			t.Errorf("c=%v T=%v: region %s, want a two-phase region", p.c, p.T, res.Region)
		}
		if math.Abs(res.TotalFraction()-1) > 1e-9 {
			t.Errorf("c=%v T=%v: mass fractions sum = %.9f, want 1", p.c, p.T, res.TotalFraction())
		}
		for _, ph := range res.Phases {
			if ph.MassFraction < 0 || ph.MassFraction > 1 {
				t.Errorf("c=%v T=%v: phase %s fraction %.6f outside [0,1]", p.c, p.T, ph.Name, ph.MassFraction)
			}
		}
	}
}

func TestEndmemberPureLiquid(t *testing.T) {
	d := testDiagram(t)
	// Pure A above its melting point must be fully liquid.
	res, err := SolvePoint(d, 0, phasediag.TAlphaMelting+10)
	if err != nil {
		t.Fatalf("pure A hot: %v", err)
	}
	if res.Region != RegionLiquid {
		t.Errorf("pure A hot: region %s, want %s", res.Region, RegionLiquid)
	}
	if w, _ := res.FractionOf("liquid"); math.Abs(w-1) > 1e-9 {
		t.Errorf("pure A hot: liquid fraction %.6f, want 1", w)
	}
	// Pure A below its melting point must be single alpha.
	res, err = SolvePoint(d, 0, phasediag.TAlphaMelting-10)
	if err != nil {
		t.Fatalf("pure A solid: %v", err)
	}
	if res.Region != RegionAlpha {
		t.Errorf("pure A solid: region %s, want %s", res.Region, RegionAlpha)
	}
	if w, _ := res.FractionOf("alpha"); math.Abs(w-1) > 1e-9 {
		t.Errorf("pure A solid: alpha fraction %.6f, want 1", w)
	}
	// Pure B above its melting point must be fully liquid.
	res, err = SolvePoint(d, 1, phasediag.TBetaMelting+10)
	if err != nil {
		t.Fatalf("pure B hot: %v", err)
	}
	if res.Region != RegionLiquid {
		t.Errorf("pure B hot: region %s, want %s", res.Region, RegionLiquid)
	}
	// Pure B below its melting point must be single beta.
	res, err = SolvePoint(d, 1, phasediag.TBetaMelting-10)
	if err != nil {
		t.Fatalf("pure B solid: %v", err)
	}
	if res.Region != RegionBeta {
		t.Errorf("pure B solid: region %s, want %s", res.Region, RegionBeta)
	}
	if w, _ := res.FractionOf("beta"); math.Abs(w-1) > 1e-9 {
		t.Errorf("pure B solid: beta fraction %.6f, want 1", w)
	}
}

func TestEutecticJustBelowTE(t *testing.T) {
	d := testDiagram(t)
	// c = cE just above TE must be fully liquid.
	res, err := SolvePoint(d, phasediag.CEutectic, phasediag.TEutectic+1)
	if err != nil {
		t.Fatalf("c=cE above TE: %v", err)
	}
	if res.Region != RegionLiquid {
		t.Errorf("c=cE above TE: region %s, want %s", res.Region, RegionLiquid)
	}
	// c = cE just below TE must be alpha + beta with lever fractions summing to 1.
	res, err = SolvePoint(d, phasediag.CEutectic, phasediag.TEutectic-1)
	if err != nil {
		t.Fatalf("c=cE below TE: %v", err)
	}
	if res.Region != RegionAlphaBeta {
		t.Errorf("c=cE below TE: region %s, want %s", res.Region, RegionAlphaBeta)
	}
	wantWAlpha := (phasediag.CBetaMin - phasediag.CEutectic) / (phasediag.CBetaMin - phasediag.CAlphaMax)
	wAlpha, _ := res.FractionOf("alpha")
	if math.Abs(wAlpha-wantWAlpha) > 1e-9 {
		t.Errorf("c=cE below TE: alpha fraction %.6f, want %.6f", wAlpha, wantWAlpha)
	}
	if math.Abs(res.TotalFraction()-1) > 1e-9 {
		t.Errorf("c=cE below TE: mass fractions sum %.9f, want 1", res.TotalFraction())
	}
	if res.EutecticFraction == nil || math.Abs(*res.EutecticFraction-1) > 1e-9 {
		t.Errorf("c=cE below TE: eutectic fraction %v, want 1", res.EutecticFraction)
	}
	// A hypoeutectic alloy retains less than one eutectic fraction.
	res, err = SolvePoint(d, 0.4, phasediag.TEutectic-1)
	if err != nil {
		t.Fatalf("hypoeutectic below TE: %v", err)
	}
	if res.EutecticFraction == nil || *res.EutecticFraction <= 0 || *res.EutecticFraction >= 1 {
		t.Errorf("hypoeutectic below TE: eutectic fraction %v, want in (0,1)", res.EutecticFraction)
	}
}

func TestBoundaryPhaseFractionOneZero(t *testing.T) {
	d := testDiagram(t)
	T := 480.0
	b, err := BoundariesAt(d, T)
	if err != nil {
		t.Fatalf("boundaries at T=%v: %v", T, err)
	}
	// c equal to the alpha boundary gives alpha fraction 1, liquid 0.
	res, err := SolvePoint(d, b.AlphaComp, T)
	if err != nil {
		t.Fatalf("SolvePoint at alpha boundary: %v", err)
	}
	if res.Region != RegionAlpha {
		t.Errorf("at alpha boundary region = %s, want %s", res.Region, RegionAlpha)
	}
	// c equal to the liquidus-left boundary gives fully liquid.
	res, err = SolvePoint(d, b.LiquidLeft, T)
	if err != nil {
		t.Fatalf("SolvePoint at liquidus boundary: %v", err)
	}
	if res.Region != RegionLiquid {
		t.Errorf("at liquidus boundary region = %s, want %s", res.Region, RegionLiquid)
	}
	if w, _ := res.FractionOf("liquid"); math.Abs(w-1) > 1e-9 {
		t.Errorf("at liquidus boundary liquid fraction %.6f, want 1", w)
	}
}

func TestInvalidInputsRejected(t *testing.T) {
	d := testDiagram(t)
	if _, err := SolvePoint(d, -0.1, 500); err == nil {
		t.Error("negative composition must be rejected")
	}
	if _, err := SolvePoint(d, 1.2, 500); err == nil {
		t.Error("composition above 1 must be rejected")
	}
	if _, err := SolvePoint(d, 0.5, -5); err == nil {
		t.Error("temperature below absolute zero must be rejected")
	}
	if _, err := SolveScan(d, 0.5, 500, 500, 10); err == nil {
		t.Error("scan with tmax <= tmin must be rejected")
	}
	if _, err := SolveScan(d, 0.5, 400, 600, 1); err == nil {
		t.Error("scan with fewer than 2 points must be rejected")
	}
	if _, err := SolveScan(d, 0.5, 400, 600, 3000); err == nil {
		t.Error("scan with too many points must be rejected")
	}
}

func TestScanLiquidFractionCurve(t *testing.T) {
	d := testDiagram(t)
	c := 0.45
	res, err := SolveScan(d, c, 300, 700, 81)
	if err != nil {
		t.Fatalf("SolveScan: %v", err)
	}
	if len(res.Points) != 81 {
		t.Fatalf("scan produced %d points, want 81", len(res.Points))
	}
	// Cold end is fully solid, hot end is fully liquid.
	if math.Abs(res.Points[0].LiquidFraction) > 1e-9 {
		t.Errorf("cold end fL = %.6f, want 0", res.Points[0].LiquidFraction)
	}
	last := res.Points[len(res.Points)-1]
	if math.Abs(last.LiquidFraction-1) > 1e-9 {
		t.Errorf("hot end fL = %.6f, want 1", last.LiquidFraction)
	}
	// fL must stay inside [0,1] and be non-decreasing in temperature.
	sawIntermediate := false
	for i, p := range res.Points {
		if p.LiquidFraction < 0 || p.LiquidFraction > 1 {
			t.Errorf("point %d fL = %.6f outside [0,1]", i, p.LiquidFraction)
		}
		if p.LiquidFraction > 0 && p.LiquidFraction < 1 {
			sawIntermediate = true
		}
		if i > 0 && p.LiquidFraction+1e-9 < res.Points[i-1].LiquidFraction {
			t.Errorf("fL not non-decreasing at point %d: %.6f after %.6f", i, p.LiquidFraction, res.Points[i-1].LiquidFraction)
		}
	}
	if !sawIntermediate {
		t.Error("scan curve has no two-phase points, want a real lever curve")
	}
}

func TestScanClosedFormMatch(t *testing.T) {
	d := testDiagram(t)
	c, T := 0.45, 470.0
	cAlpha := phasediag.CAlphaMax * (phasediag.TAlphaMelting - T) / (phasediag.TAlphaMelting - phasediag.TEutectic)
	cL := phasediag.CEutectic * (phasediag.TAlphaMelting - T) / (phasediag.TAlphaMelting - phasediag.TEutectic)
	want := (c - cAlpha) / (cL - cAlpha)

	res, err := SolveScan(d, c, T-0.1, T+0.1, 11)
	if err != nil {
		t.Fatalf("SolveScan: %v", err)
	}
	// The point nearest T must match the closed form within scan spacing tolerance.
	closest := res.Points[0]
	for _, p := range res.Points {
		if math.Abs(p.Temperature-T) < math.Abs(closest.Temperature-T) {
			closest = p
		}
	}
	if math.Abs(closest.LiquidFraction-want) > 2e-3 {
		t.Errorf("scan fL at T=%v = %.6f, closed form want %.6f", closest.Temperature, closest.LiquidFraction, want)
	}
}
