package web

import (
	"net/http"

	"lever-phase/internal/lever"
)

// handleScan 处理 POST /api/scan：固定成分 c，在 [tmin, tmax] 上
// 扫描 n 个温度点，返回 fL(T) 点列。点列全部来自后端求解。
func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req scanRequest
	if err := decodeBody(r, &req); err != nil {
		badRequest(w, err)
		return
	}
	n := req.N
	if n == 0 {
		n = lever.DefaultScanPoints
	}
	if err := lever.ValidateScanInput(req.C, req.TMin, req.TMax, n); err != nil {
		badRequest(w, err)
		return
	}
	res, err := lever.SolveScan(s.diagram, req.C, req.TMin, req.TMax, n)
	if err != nil {
		badRequest(w, err)
		return
	}
	writeJSON(w, http.StatusOK, scanResponse{
		Composition: res.Composition,
		TEutecticK:  s.diagram.EutecticT,
		Points:      res.Points,
	})
}
