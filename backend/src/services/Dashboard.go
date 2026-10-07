package services

import (
	"groundTurn/src/constants"
	"groundTurn/src/types"
)

type DashboardService struct{ viewService *ViewService }

func NewDashboardService(vs *ViewService) *DashboardService {
	return &DashboardService{viewService: vs}
}

// Build 看板汇总。超时任务与任务列表同源（ViewService.LoadBundle），口径天然一致。
func (s *DashboardService) Build() (*types.DashboardView, error) {
	bundle, err := s.viewService.LoadBundle()
	if err != nil {
		return nil, err
	}
	out := &types.DashboardView{
		Turnarounds: bundle.Turnarounds,
		TaskTotal:   len(bundle.Tasks),
	}
	for i := range bundle.Turnarounds {
		t := &bundle.Turnarounds[i]
		out.TurnaroundTotal++
		if t.TurnaroundStatus != constants.TurnaroundDeparted {
			out.ActiveTurnaround++
		}
		out.OpenDelayMinutes += t.OpenDelayMinutes
		out.OpenDelayEvents += t.OpenDelayCount
	}
	for i := range bundle.Tasks {
		t := &bundle.Tasks[i]
		if t.Overdue {
			out.OverdueTaskCount++
			out.OverdueTasks = append(out.OverdueTasks, *t)
		}
		if t.Status == constants.TaskStatusCompleted {
			out.TaskCompleted++
		}
		if t.AcceptedAt != nil {
			out.TaskAccepted++
		}
	}
	return out, nil
}
