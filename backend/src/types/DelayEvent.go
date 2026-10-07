package types

import "time"

// DelayEventView 延误事件响应对象。ResolvedAt 为 nil 即未关闭，仍在顺延任务截止时间。
type DelayEventView struct {
	ID                 uint       `json:"id"`
	TurnaroundID       uint       `json:"turnaround_id"`
	FlightNo           string     `json:"flight_no"`
	DelayType          string     `json:"delay_type"`
	Minutes            int        `json:"minutes"`
	RootCause          string     `json:"root_cause"`
	ResponsibilityTeam string     `json:"responsibility_team"`
	CreatedAt          time.Time  `json:"created_at"`
	ResolvedAt         *time.Time `json:"resolved_at"`
	Closed             bool       `json:"closed"`
	// OpenDelayMinutes 该航班当前所有未关闭延误累计分钟（本事件关闭与否都会带出当前口径）。
	OpenDelayMinutes int `json:"open_delay_minutes"`
}

// DelayEventRegisterRequest 登记延误。登记后立即顺延该航班全部未完成任务的生效截止（不改 deadline 列）。
type DelayEventRegisterRequest struct {
	TurnaroundID       uint   `json:"turnaround_id" binding:"required"`
	DelayType          string `json:"delay_type" binding:"required"`
	Minutes            int    `json:"minutes" binding:"required"`
	RootCause          string `json:"root_cause"`
	ResponsibilityTeam string `json:"responsibility_team"`
}
