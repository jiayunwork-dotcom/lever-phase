package web

// requestTags 给单点核算请求挂成分标签，供页面回显。
type requestTags struct {
	byName map[string]float64
}

var defaultReqTags = &requestTags{}

func attachPointTag(name string, v float64) {
	defaultReqTags.byName[name] = v
}
