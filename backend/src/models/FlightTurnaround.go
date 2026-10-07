package models

import "time"

// FlightTurnaround 航班过站。DelayEvents 与 GroundTask 决定任务的生效截止时间。
type FlightTurnaround struct {
	ID               uint         `gorm:"primaryKey" json:"id"`
	FlightNo         string       `gorm:"size:32;not null" json:"flight_no"`
	AircraftReg      string       `gorm:"size:32" json:"aircraft_reg"`
	StandNo          string       `gorm:"size:16" json:"stand_no"`
	ArrivalTime      time.Time    `json:"arrival_time"`
	DepartureTime    time.Time    `json:"departure_time"`
	TurnaroundStatus string       `gorm:"size:32;not null;index" json:"turnaround_status"`
	DelayReason      string       `gorm:"size:255" json:"delay_reason"`
	Tasks            []GroundTask `gorm:"foreignKey:TurnaroundID" json:"-"`
	DelayEvents      []DelayEvent `gorm:"foreignKey:TurnaroundID" json:"-"`
}

func (FlightTurnaround) TableName() string { return "flight_turnaround" }
