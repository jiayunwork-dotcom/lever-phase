package web

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

const alloy60Example = `{
  "label": "Pb-60Sn alloy in the alpha + liquid two-phase region",
  "c": 0.6,
  "t": 458.0
}`

func testHandler(t *testing.T) http.Handler {
	t.Helper()
	webFS := fstest.MapFS{
		"web/index.html": &fstest.MapFile{Data: []byte("<h1>lever-phase</h1>")},
		"web/app.js":     &fstest.MapFile{Data: []byte("// test stub")},
	}
	h, err := NewServer(Assets{
		WebFS: webFS,
		Examples: map[string][]byte{
			"alloy-60": []byte(alloy60Example),
		},
	})
	if err != nil {
		t.Fatalf("NewServer: %v", err)
	}
	return h
}

func post(t *testing.T, handler http.Handler, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func get(t *testing.T, handler http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func TestPointEndpointOK(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/point", `{"c":0.45,"t":475}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp pointResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Region != "alpha_liquid" {
		t.Errorf("region = %s, want alpha_liquid", resp.Region)
	}
	if len(resp.Phases) != 2 {
		t.Fatalf("len(phases) = %d, want 2", len(resp.Phases))
	}
	if math.Abs(resp.TotalFraction-1) > 1e-9 {
		t.Errorf("total fraction = %.9f, want 1", resp.TotalFraction)
	}
	if resp.LiquidFraction <= 0 || resp.LiquidFraction >= 1 {
		t.Errorf("liquid fraction = %.6f, want in (0,1)", resp.LiquidFraction)
	}
}

func TestPointEndpointNegativeComposition(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/point", `{"c":-0.1,"t":500}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "composition") {
		t.Errorf("error = %q, want it to mention composition", resp.Error)
	}
}

func TestPointEndpointSubZeroTemperature(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/point", `{"c":0.5,"t":-5}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "absolute zero") {
		t.Errorf("error = %q, want it to mention absolute zero", resp.Error)
	}
}

func TestScanEndpointOK(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/scan", `{"c":0.45,"tmin":300,"tmax":700,"n":81}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp scanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Points) != 81 {
		t.Errorf("len(points) = %d, want 81", len(resp.Points))
	}
	if math.Abs(resp.Points[0].LiquidFraction) > 1e-9 {
		t.Errorf("cold end fL = %.6f, want 0", resp.Points[0].LiquidFraction)
	}
	last := resp.Points[len(resp.Points)-1]
	if math.Abs(last.LiquidFraction-1) > 1e-9 {
		t.Errorf("hot end fL = %.6f, want 1", last.LiquidFraction)
	}
	if math.Abs(resp.TEutecticK-456) > 1e-9 {
		t.Errorf("teutectic_k = %v, want 456", resp.TEutecticK)
	}
}

func TestScanEndpointInvalidRange(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/scan", `{"c":0.5,"tmin":500,"tmax":500,"n":21}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "tmax") {
		t.Errorf("error = %q, want it to mention tmax", resp.Error)
	}
}

func TestScanEndpointSubZeroRange(t *testing.T) {
	handler := testHandler(t)
	rec := post(t, handler, "/api/scan", `{"c":0.5,"tmin":-10,"tmax":700,"n":21}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resp.Error, "absolute zero") {
		t.Errorf("error = %q, want it to mention absolute zero", resp.Error)
	}
}

func TestExamplesEndpoint(t *testing.T) {
	handler := testHandler(t)
	rec := get(t, handler, "/api/examples")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		Alloy60 json.RawMessage `json:"alloy-60"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp.Alloy60) == 0 {
		t.Fatal("examples payload missing alloy-60")
	}
	var ex pointRequest
	if err := json.Unmarshal(resp.Alloy60, &ex); err != nil {
		t.Fatal(err)
	}
	if math.Abs(ex.C-0.6) > 1e-9 || math.Abs(ex.T-458) > 1e-9 {
		t.Errorf("example parsed as c=%v t=%v, want c=0.6 t=458", ex.C, ex.T)
	}
}

func TestStaticIndexServed(t *testing.T) {
	handler := testHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 for index", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "lever-phase") {
		t.Errorf("index page body does not mention lever-phase")
	}
}
