package models

import "time"

// ResourceBooking 资源预约，连接航班、任务和资源。
type ResourceBooking struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ResourceID     uint      `gorm:"not null;index" json:"resource_id"`
	TurnaroundID   uint      `gorm:"not null;index" json:"turnaround_id"`
	TaskID         *uint     `json:"task_id"`
	StartTime      time.Time `gorm:"not null" json:"start_time"`
	EndTime        time.Time `gorm:"not null" json:"end_time"`
	BookingStatus  string    `gorm:"size:32;not null" json:"booking_status"`
	ConflictReason string    `gorm:"size:255" json:"conflict_reason"`
}

func (ResourceBooking) TableName() string { return "resource_booking" }
