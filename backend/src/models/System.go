package models

import "time"

// User 登录账号，承载 RBAC 角色。
type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"size:64;uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"size:255;not null" json:"-"`
	Role         string `gorm:"size:32;not null" json:"role"`
	TeamCode     string `gorm:"size:32" json:"team_code"`
	DisplayName  string `gorm:"size:64" json:"display_name"`
}

func (User) TableName() string { return "app_user" }

// AuditLog 操作日志（派工、资源预约、延误归因等写操作）。
type AuditLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Actor      string    `gorm:"size:64" json:"actor"`
	Action     string    `gorm:"size:64;not null" json:"action"`
	TargetType string    `gorm:"size:32;not null" json:"target_type"`
	TargetID   string    `gorm:"size:64" json:"target_id"`
	Detail     string    `gorm:"size:512" json:"detail"`
	CreatedAt  time.Time `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_log" }
