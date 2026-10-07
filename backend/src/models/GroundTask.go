package models

import "time"

// GroundTask 地勤任务。
// Deadline 是原计划截止时间，登记延误永不改写该列；
// 顺延结果由服务层按航班上未关闭的 DelayEvent 另算为 effective_deadline。
type GroundTask struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	TurnaroundID uint       `gorm:"not null;index" json:"turnaround_id"`
	TaskType     string     `gorm:"size:32;not null" json:"task_type"`
	TeamID       string     `gorm:"size:32" json:"team_id"`
	PlannedStart *time.Time `json:"planned_start"`
	// Deadline 原计划截止时间（保留不变，是顺延的基准）。
	Deadline time.Time `gorm:"not null" json:"deadline"`
	// AcceptedAt 签收时间；未签收为 nil。
	AcceptedAt   *time.Time `json:"accepted_at"`
	ActualFinish *time.Time `json:"actual_finish"`
	Status       string     `gorm:"size:32;not null;index" json:"status"`
	BlockerNote  string     `gorm:"size:255" json:"blocker_note"`
}

func (GroundTask) TableName() string { return "ground_task" }
