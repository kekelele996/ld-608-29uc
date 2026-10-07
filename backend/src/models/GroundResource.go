package models

import "time"

// GroundResource 保障资源台账。
type GroundResource struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	ResourceCode       string     `gorm:"size:32;not null" json:"resource_code"`
	ResourceType       string     `gorm:"size:32;not null" json:"resource_type"`
	Location           string     `gorm:"size:64" json:"location"`
	AvailabilityStatus string     `gorm:"size:32;not null;index" json:"availability_status"`
	MaintenanceDueAt   *time.Time `json:"maintenance_due_at"`
	OwnerTeam          string     `gorm:"size:32" json:"owner_team"`
}

func (GroundResource) TableName() string { return "ground_resource" }
