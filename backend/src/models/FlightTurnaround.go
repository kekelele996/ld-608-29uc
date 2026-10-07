package models

type FlightTurnaround struct {
	ID               int    `json:"id"`
	FlightNo         string `json:"flight_no"`
	AircraftReg      string `json:"aircraft_reg"`
	StandNo          string `json:"stand_no"`
	ArrivalTime      string `json:"arrival_time"`
	DepartureTime    string `json:"departure_time"`
	TurnaroundStatus string `json:"turnaround_status"`
	DelayReason      string `json:"delay_reason"`
}
