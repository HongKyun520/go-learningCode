package web

type Result struct {
	// 业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Date any    `json:"date"`
}
