package types

// CreateDelayEventRequest 登记延误请求；新事件一律未关闭（resolved_at 为空）。
type CreateDelayEventRequest struct {
	TurnaroundID       int    `json:"turnaround_id"`
	DelayType          string `json:"delay_type"`
	Minutes            int    `json:"minutes"`
	RootCause          string `json:"root_cause"`
	ResponsibilityTeam string `json:"responsibility_team"`
}

// CloseDelayEventRequest 关闭延误请求；closed_at 缺省时由服务端取当前时间。
type CloseDelayEventRequest struct {
	ClosedAt string `json:"closed_at"`
}
