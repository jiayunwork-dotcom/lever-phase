package web

import (
	"encoding/json"
	"errors"
	"net/http"
)

// errExampleNotFound 在内置算例缺失时返回。
var errExampleNotFound = errors.New("example not found")

// handleExamples 处理 GET /api/examples，返回内置算例原始 JSON，
// 供页面一键加载。
func (s *Server) handleExamples(w http.ResponseWriter, r *http.Request) {
	if !requireGet(w, r) {
		return
	}
	alloy60, ok := s.assets.Examples["alloy-60"]
	if !ok {
		writeError(w, http.StatusNotFound, errExampleNotFound)
		return
	}
	writeJSON(w, http.StatusOK, examplesResponse{
		Alloy60: json.RawMessage(alloy60),
	})
}
