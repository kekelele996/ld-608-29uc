package types

import "time"

// FlightTurnaroundView 航班过站响应对象，含未关闭延误口径。
type FlightTurnaroundView struct {
	ID               uint             `json:"id"`
	FlightNo         string           `json:"flight_no"`
	AircraftReg      string           `json:"aircraft_reg"`
	StandNo          string           `json:"stand_no"`
	ArrivalTime      time.Time        `json:"arrival_time"`
	DepartureTime    time.Time        `json:"departure_time"`
	TurnaroundStatus string           `json:"turnaround_status"`
	StatusText       string           `json:"status_text"`
	DelayReason      string           `json:"delay_reason"`
	OpenDelayMinutes int              `json:"open_delay_minutes"`
	OpenDelayCount   int              `json:"open_delay_count"`
	TaskTotal        int              `json:"task_total"`
	TaskCompleted    int              `json:"task_completed"`
	TaskOverdue      int              `json:"task_overdue"`
	Tasks            []GroundTaskView `json:"tasks,omitempty"`
}

// FlightTurnaroundRegisterRequest 过站登记。
type FlightTurnaroundRegisterRequest struct {
	FlightNo      string    `json:"flight_no" binding:"required"`
	AircraftReg   string    `json:"aircraft_reg"`
	StandNo       string    `json:"stand_no"`
	ArrivalTime   time.Time `json:"arrival_time" binding:"required"`
	DepartureTime time.Time `json:"departure_time" binding:"required"`
}

// FlightTurnaroundStatusRequest 状态推进/放行。
type FlightTurnaroundStatusRequest struct {
	Status      string `json:"status" binding:"required"`
	DelayReason string `json:"delay_reason"`
}
