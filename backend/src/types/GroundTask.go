package types

import "time"

// GroundTaskView 任务响应对象。
// Deadline 始终是原计划截止时间；EffectiveDeadline 是按未关闭延误顺延后当前生效的截止时间。
// 超时判定一律以 EffectiveDeadline 为准（看板超时卡片与任务列表共用此口径）。
type GroundTaskView struct {
	ID                uint       `json:"id"`
	TurnaroundID      uint       `json:"turnaround_id"`
	FlightNo          string     `json:"flight_no"`
	TaskType          string     `json:"task_type"`
	TaskTypeText      string     `json:"task_type_text"`
	TeamID            string     `json:"team_id"`
	PlannedStart      *time.Time `json:"planned_start"`
	Deadline          time.Time  `json:"deadline"`
	EffectiveDeadline time.Time  `json:"effective_deadline"`
	OpenDelayMinutes  int        `json:"open_delay_minutes"`
	OpenDelayCount    int        `json:"open_delay_count"`
	AcceptedAt        *time.Time `json:"accepted_at"`
	ActualFinish      *time.Time `json:"actual_finish"`
	Status            string     `json:"status"`
	StatusText        string     `json:"status_text"`
	BlockerNote       string     `json:"blocker_note"`
	Overdue           bool       `json:"overdue"`
	OverdueMinutes    int        `json:"overdue_minutes"`
	CompareBasis      string     `json:"compare_basis"` // ACTUAL_FINISH / NOW
}

// GroundTaskDispatchRequest 派工（新建任务）。
type GroundTaskDispatchRequest struct {
	TurnaroundID uint       `json:"turnaround_id" binding:"required"`
	TaskType     string     `json:"task_type" binding:"required"`
	TeamID       string     `json:"team_id"`
	PlannedStart *time.Time `json:"planned_start"`
	Deadline     time.Time  `json:"deadline" binding:"required"`
}

// GroundTaskBlockRequest 阻塞登记。
type GroundTaskBlockRequest struct {
	BlockerNote string `json:"blocker_note" binding:"required"`
}
