package models

import "time"

// DelayEvent 延误事件。ResolvedAt 为 nil 表示尚未关闭，
// 只有未关闭事件的 Minutes 会计入任务截止时间的顺延。
type DelayEvent struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	TurnaroundID       uint      `gorm:"not null;index" json:"turnaround_id"`
	DelayType          string    `gorm:"size:32;not null" json:"delay_type"`
	Minutes            int       `gorm:"not null" json:"minutes"`
	RootCause          string    `gorm:"size:255" json:"root_cause"`
	ResponsibilityTeam string    `gorm:"size:32" json:"responsibility_team"`
	CreatedAt          time.Time `json:"created_at"`
	// ResolvedAt 关闭时间；NULL=未关闭（仍在顺延），非 NULL=已关闭（不再顺延）。
	ResolvedAt *time.Time `json:"resolved_at"`
}

func (DelayEvent) TableName() string { return "delay_event" }
