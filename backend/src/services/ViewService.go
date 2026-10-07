package services

import (
	"time"

	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// TaskViewBundle 供看板与任务列表复用同一份装配结果，杜绝两套口径。
type TaskViewBundle struct {
	Tasks        []types.GroundTaskView
	Turnarounds  []types.FlightTurnaroundView
	FlightByID   map[uint]models.FlightTurnaround
	SummaryByTAR map[uint]utils.OpenDelaySummary
}

// ViewService 汇总 航班+任务+延误 三张表，统一计算顺延截止与超时。
type ViewService struct {
	turnaroundRepo *repositories.FlightTurnaroundRepository
	taskRepo       *repositories.GroundTaskRepository
	delayRepo      *repositories.DelayEventRepository
}

func NewViewService(
	tr *repositories.FlightTurnaroundRepository,
	kr *repositories.GroundTaskRepository,
	dr *repositories.DelayEventRepository,
) *ViewService {
	return &ViewService{turnaroundRepo: tr, taskRepo: kr, delayRepo: dr}
}

// LoadBundle 一次性装配：未关闭延误汇总 -> 生效截止 -> 超时判定。
func (s *ViewService) LoadBundle() (*TaskViewBundle, error) {
	now := time.Now().UTC()

	flights, err := s.turnaroundRepo.List()
	if err != nil {
		return nil, err
	}
	tasks, err := s.taskRepo.List()
	if err != nil {
		return nil, err
	}
	flightIDs := make([]uint, 0, len(flights))
	flightByID := map[uint]models.FlightTurnaround{}
	for i := range flights {
		flightIDs = append(flightIDs, flights[i].ID)
		flightByID[flights[i].ID] = flights[i]
	}
	delayMap, err := s.delayRepo.MapByTurnaround(flightIDs)
	if err != nil {
		return nil, err
	}

	// 每个航班只算一次未关闭延误汇总（关闭过的 resolved_at 非空，被剔除）。
	summaryByTAR := map[uint]utils.OpenDelaySummary{}
	for _, id := range flightIDs {
		summaryByTAR[id] = utils.SummarizeOpenDelays(delayMap[id])
	}

	taskViews := make([]types.GroundTaskView, 0, len(tasks))
	viewsByTAR := map[uint][]types.GroundTaskView{}
	for i := range tasks {
		t := tasks[i]
		flightNo := flightByID[t.TurnaroundID].FlightNo
		view := constructors.BuildGroundTaskView(t, flightNo, summaryByTAR[t.TurnaroundID], now)
		taskViews = append(taskViews, view)
		viewsByTAR[t.TurnaroundID] = append(viewsByTAR[t.TurnaroundID], view)
	}

	turnaroundViews := make([]types.FlightTurnaroundView, 0, len(flights))
	for i := range flights {
		turnaroundViews = append(turnaroundViews, constructors.BuildFlightTurnaroundView(
			flights[i], viewsByTAR[flights[i].ID], summaryByTAR[flights[i].ID],
		))
	}

	return &TaskViewBundle{
		Tasks:        taskViews,
		Turnarounds:  turnaroundViews,
		FlightByID:   flightByID,
		SummaryByTAR: summaryByTAR,
	}, nil
}
