package web

import (
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// staticHandler 从嵌入的 fs.FS 中取 web 子目录提供静态文件；
// 找不到路径时回退到 index.html。
func staticHandler(webFS fs.FS) http.Handler {
	if webFS == nil {
		return http.NotFoundHandler()
	}
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		return http.NotFoundHandler()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := path.Clean(strings.TrimPrefix(r.URL.Path, "/"))
		if clean == "." || clean == "/" || clean == "" {
			clean = "index.html"
		}
		if _, err := fs.Stat(sub, clean); err != nil {
			clean = "index.html"
		}
		serveFile(w, r, sub, clean)
	})
}

// serveFile 打开文件并按内容类型返回。
func serveFile(w http.ResponseWriter, r *http.Request, root fs.FS, name string) {
	f, err := root.Open(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}
	rs, ok := f.(io.ReadSeeker)
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, info.Name(), info.ModTime(), rs)
}
