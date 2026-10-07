package types

// SignOffGroundTaskRequest 签收请求；signed_at 缺省时由服务端取当前时间。
type SignOffGroundTaskRequest struct {
	SignedAt string `json:"signed_at"`
}
