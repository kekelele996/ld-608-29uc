package constructors

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// NewGroundTask 从派工请求构造任务：初始待签收，截止时间取原计划。
func NewGroundTask(req types.GroundTaskDispatchRequest) models.GroundTask {
	return models.GroundTask{
		TurnaroundID: req.TurnaroundID,
		TaskType:     req.TaskType,
		TeamID:       req.TeamID,
		PlannedStart: req.PlannedStart,
		Deadline:     req.Deadline.UTC(),
		Status:       constants.TaskStatusPending,
	}
}

// BuildGroundTaskView 组装任务响应，落地统一口径：
// 保留原计划 deadline，按航班未关闭延误另算 effective_deadline，并据此判超时。
func BuildGroundTaskView(t models.GroundTask, flightNo string, summary utils.OpenDelaySummary, now time.Time) types.GroundTaskView {
	effective := utils.EffectiveDeadline(t.Deadline, summary)
	// 完成过的任务按 actual_finish 比生效截止；未完成（含仅签收）按当前时间，actual_finish 为 nil。
	verdict := utils.EvaluateOverdue(t.Status, effective, t.ActualFinish, now)
	return types.GroundTaskView{
		ID:                t.ID,
		TurnaroundID:      t.TurnaroundID,
		FlightNo:          flightNo,
		TaskType:          t.TaskType,
		TaskTypeText:      constants.StatusText["GroundTaskType"][t.TaskType],
		TeamID:            t.TeamID,
		PlannedStart:      t.PlannedStart,
		Deadline:          t.Deadline,
		EffectiveDeadline: effective,
		OpenDelayMinutes:  summary.OpenMinutes,
		OpenDelayCount:    summary.OpenCount,
		AcceptedAt:        t.AcceptedAt,
		ActualFinish:      t.ActualFinish,
		Status:            t.Status,
		StatusText:        constants.StatusText["GroundTaskStatus"][t.Status],
		BlockerNote:       t.BlockerNote,
		Overdue:           verdict.Overdue,
		OverdueMinutes:    verdict.OverdueMinutes,
		CompareBasis:      verdict.CompareBasis,
	}
}
