package types

import "time"

type ResourceBookingView struct {
	ID             uint      `json:"id"`
	ResourceID     uint      `json:"resource_id"`
	ResourceCode   string    `json:"resource_code"`
	TurnaroundID   uint      `json:"turnaround_id"`
	TaskID         *uint     `json:"task_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	BookingStatus  string    `json:"booking_status"`
	StatusText     string    `json:"status_text"`
	ConflictReason string    `json:"conflict_reason"`
}

type ResourceBookingRequest struct {
	ResourceID   uint      `json:"resource_id" binding:"required"`
	TurnaroundID uint      `json:"turnaround_id" binding:"required"`
	TaskID       *uint     `json:"task_id"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	EndTime      time.Time `json:"end_time" binding:"required"`
}

// DashboardView 看板：超时任务卡片与任务列表共用 GroundTaskView 口径。
type DashboardView struct {
	TurnaroundTotal  int                    `json:"turnaround_total"`
	ActiveTurnaround int                    `json:"active_turnaround"`
	OpenDelayMinutes int                    `json:"open_delay_minutes"`
	OpenDelayEvents  int                    `json:"open_delay_events"`
	OverdueTaskCount int                    `json:"overdue_task_count"`
	TaskTotal        int                    `json:"task_total"`
	TaskAccepted     int                    `json:"task_accepted"`
	TaskCompleted    int                    `json:"task_completed"`
	OverdueTasks     []GroundTaskView       `json:"overdue_tasks"`
	Turnarounds      []FlightTurnaroundView `json:"turnarounds"`
}
