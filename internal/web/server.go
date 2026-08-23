// Package web 提供 lever-phase 的 HTTP 控制台：/api/point 单点核算、
// /api/scan 温度扫描、/api/examples 内置算例，以及由 Go 同进程托管
// 的静态页面。所有错误都返回带 error 字段的 JSON，页面与 API 都能
// 读到后端失败原因。
package web

import (
	"io/fs"
	"net/http"

	"lever-phase/internal/phasediag"
)

// Assets 是构造服务所需的静态资源与内置算例。
type Assets struct {
	WebFS    fs.FS
	Examples map[string][]byte
}

// Server 持有领域内核、静态资源与路由 mux。
type Server struct {
	diagram phasediag.PhaseDiagram
	assets  Assets
	mux     *http.ServeMux
}

// NewServer 注册全部路由并返回可用的 http.Handler。默认使用经典
// Pb–Sn 简化相图；静态资源缺失时仍可提供 API。
func NewServer(assets Assets) (http.Handler, error) {
	diagram := phasediag.DefaultPbSn()
	if err := diagram.Validate(); err != nil {
		return nil, err
	}
	s := &Server{diagram: diagram, assets: assets, mux: http.NewServeMux()}
	s.mux.HandleFunc("/api/point", s.handlePoint)
	s.mux.HandleFunc("/api/scan", s.handleScan)
	s.mux.HandleFunc("/api/examples", s.handleExamples)
	s.mux.Handle("/", staticHandler(assets.WebFS))
	return s.mux, nil
}

// Diagram 暴露当前使用的相图（供调试与测试）。
func (s *Server) Diagram() phasediag.PhaseDiagram { return s.diagram }
