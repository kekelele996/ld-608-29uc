package repositories

import (
	"sync"

	"groundTurn/src/models"
)

// 内存种子与 frontend/src/mocks/seedData.ts 保持一致。
var flightTurnaroundMu sync.RWMutex

var flightTurnaroundRows = []models.FlightTurnaround{
	{ID: 1, FlightNo: "CA1234", AircraftReg: "B-6688", StandNo: "12", ArrivalTime: "2026-10-07T04:50:00Z", DepartureTime: "2026-10-07T08:30:00Z", TurnaroundStatus: "IN_SERVICE", DelayReason: "加油车调度延迟"},
	{ID: 2, FlightNo: "MU5678", AircraftReg: "B-8866", StandNo: "07", ArrivalTime: "2026-10-07T06:10:00Z", DepartureTime: "2026-10-07T09:40:00Z", TurnaroundStatus: "ON_STAND", DelayReason: "行李分拣系统故障"},
	{ID: 3, FlightNo: "ZH9012", AircraftReg: "B-1234", StandNo: "21", ArrivalTime: "2026-10-07T08:20:00Z", DepartureTime: "2026-10-07T11:00:00Z", TurnaroundStatus: "DELAYED", DelayReason: "前序航班晚到"},
}

func ListFlightTurnarounds() []models.FlightTurnaround {
	flightTurnaroundMu.RLock()
	defer flightTurnaroundMu.RUnlock()
	rows := make([]models.FlightTurnaround, len(flightTurnaroundRows))
	copy(rows, flightTurnaroundRows)
	return rows
}
