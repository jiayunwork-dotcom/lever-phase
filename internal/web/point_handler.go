package web

import (
	"net/http"

	"lever-phase/internal/lever"
)

// handlePoint 处理 POST /api/point：由请求体中的 c、t 做一次杠杆
// 核算，成功返回相区与相分数，非法参数返回 400 + error JSON。
func (s *Server) handlePoint(w http.ResponseWriter, r *http.Request) {
	if !requirePost(w, r) {
		return
	}
	var req pointRequest
	if err := decodeBody(r, &req); err != nil {
		badRequest(w, err)
		return
	}
	if err := lever.ValidatePointRequest(req.C, req.T); err != nil {
		badRequest(w, err)
		return
	}
	res, err := lever.SolvePoint(s.diagram, req.C, req.T)
	if err != nil {
		badRequest(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toPointResponse(res))
}
