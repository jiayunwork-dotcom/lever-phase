package web

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

const maxBodyBytes = 1 << 20

// errorBody 是统一错误响应体。
type errorBody struct {
	Error string `json:"error"`
}

// writeJSON 写出 JSON 响应。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

// writeError 把 error 写为带 error 字段的 JSON。
func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, errorBody{Error: err.Error()})
}

// badRequest 以 400 写出错误体。
func badRequest(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, err)
}

// decodeBody 读取并解析 JSON 请求体，带大小上限；解析失败返回
// 可读的错误文案而不是吞掉。
func decodeBody(r *http.Request, v any) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	defer r.Body.Close()
	if r.ContentLength > maxBodyBytes {
		return errors.New("request body too large")
	}
	data, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		return errors.New("read body: " + err.Error())
	}
	if int64(len(data)) > maxBodyBytes {
		return errors.New("request body too large")
	}
	if len(data) == 0 {
		return errors.New("empty request body")
	}
	if err := json.Unmarshal(data, v); err != nil {
		return errors.New("invalid JSON: " + err.Error())
	}
	return nil
}

// requirePost 校验 POST 方法，否则返回 405 并写错误体。
func requirePost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodPost {
		return true
	}
	writeJSON(w, http.StatusMethodNotAllowed, errorBody{Error: "method " + r.Method + " not allowed on this endpoint"})
	return false
}

// requireGet 校验 GET 方法。
func requireGet(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet {
		return true
	}
	writeJSON(w, http.StatusMethodNotAllowed, errorBody{Error: "method " + r.Method + " not allowed on this endpoint"})
	return false
}
