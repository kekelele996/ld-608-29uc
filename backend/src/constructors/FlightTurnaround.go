package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// NewFlightTurnaround 过站登记构造，初始状态 ARRIVING。
func NewFlightTurnaround(req types.FlightTurnaroundRegisterRequest) models.FlightTurnaround {
	return models.FlightTurnaround{
		FlightNo:         req.FlightNo,
		AircraftReg:      req.AircraftReg,
		StandNo:          req.StandNo,
		ArrivalTime:      req.ArrivalTime.UTC(),
		DepartureTime:    req.DepartureTime.UTC(),
		TurnaroundStatus: constants.TurnaroundArriving,
	}
}

// BuildFlightTurnaroundView 组装航班响应，任务列表使用同一套顺延/超时口径。
func BuildFlightTurnaroundView(f models.FlightTurnaround, taskViews []types.GroundTaskView, summary utils.OpenDelaySummary) types.FlightTurnaroundView {
	view := types.FlightTurnaroundView{
		ID:               f.ID,
		FlightNo:         f.FlightNo,
		AircraftReg:      f.AircraftReg,
		StandNo:          f.StandNo,
		ArrivalTime:      f.ArrivalTime,
		DepartureTime:    f.DepartureTime,
		TurnaroundStatus: f.TurnaroundStatus,
		StatusText:       constants.StatusText["TurnaroundStatus"][f.TurnaroundStatus],
		DelayReason:      f.DelayReason,
		OpenDelayMinutes: summary.OpenMinutes,
		OpenDelayCount:   summary.OpenCount,
		TaskTotal:        len(taskViews),
		Tasks:            taskViews,
	}
	for _, t := range taskViews {
		if t.Status == constants.TaskStatusCompleted {
			view.TaskCompleted++
		}
		if t.Overdue {
			view.TaskOverdue++
		}
	}
	return view
}
