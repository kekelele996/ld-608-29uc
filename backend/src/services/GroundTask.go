package services

import (
	"errors"
	"fmt"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/middlewares"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GroundTaskService 派工、签收、完成、阻塞。
// 关键约束：任何动作都不改写 deadline 列；超时一律按顺延后的 effective_deadline。
type GroundTaskService struct {
	taskRepo       *repositories.GroundTaskRepository
	turnaroundRepo *repositories.FlightTurnaroundRepository
	delayRepo      *repositories.DelayEventRepository
}

func NewGroundTaskService(
	kr *repositories.GroundTaskRepository,
	tr *repositories.FlightTurnaroundRepository,
	dr *repositories.DelayEventRepository,
) *GroundTaskService {
	return &GroundTaskService{taskRepo: kr, turnaroundRepo: tr, delayRepo: dr}
}

func (s *GroundTaskService) validateEnum(field string, allowed []string, value string) error {
	for _, v := range allowed {
		if v == value {
			return nil
		}
	}
	return utils.EnumError(field, value)
}

// summaryOf 取航班上未关闭延误汇总（重算生效截止用）。
func (s *GroundTaskService) summaryOf(turnaroundID uint) (utils.OpenDelaySummary, error) {
	events, err := s.delayRepo.MapByTurnaround([]uint{turnaroundID})
	if err != nil {
		return utils.OpenDelaySummary{}, err
	}
	return utils.SummarizeOpenDelays(events[turnaroundID]), nil
}

func (s *GroundTaskService) loadTask(id uint) (*models.GroundTask, error) {
	t, err := s.taskRepo.Get(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.NotFoundError(constants.TargetTask, id)
	}
	return t, err
}

// Dispatch 派工：写入原计划截止，不做任何延误改写。
func (s *GroundTaskService) Dispatch(c *gin.Context, req types.GroundTaskDispatchRequest) (*types.GroundTaskView, error) {
	if err := s.validateEnum("task_type", constants.GroundTaskType, req.TaskType); err != nil {
		return nil, err
	}
	if _, err := s.turnaroundRepo.Get(req.TurnaroundID); err != nil {
		return nil, utils.NotFoundError(constants.TargetTurnaround, req.TurnaroundID)
	}
	if !req.Deadline.IsZero() && req.PlannedStart != nil && req.Deadline.Before(*req.PlannedStart) {
		return nil, utils.ValidationError(fmt.Sprintf(constants.ValidationFailedMessage, "deadline must be after planned_start"))
	}
	task := constructors.NewGroundTask(req)
	if err := s.taskRepo.Create(&task); err != nil {
		return nil, err
	}
	summary, _ := s.summaryOf(task.TurnaroundID)
	view := constructors.BuildGroundTaskView(task, "", summary, time.Now().UTC())
	flight := s.flightNo(task.TurnaroundID)
	view.FlightNo = flight

	middlewares.QueueAudit(c, "TASK_DISPATCH", constants.TargetTask, fmt.Sprint(task.ID),
		utils.RenderLog(constants.LogTemplates["TASK_DISPATCH"], map[string]interface{}{
			"TaskID": task.ID, "TaskType": task.TaskType, "TeamID": task.TeamID,
			"Deadline": task.Deadline.Format(time.RFC3339),
		}))
	return &view, nil
}

func (s *GroundTaskService) flightNo(id uint) string {
	if f, err := s.turnaroundRepo.Get(id); err == nil {
		return f.FlightNo
	}
	return ""
}

// Accept 签收：记录 accepted_at，状态置 ACCEPTED；签收动作不改截止时间。
func (s *GroundTaskService) Accept(c *gin.Context, id uint) (*types.GroundTaskView, error) {
	task, err := s.loadTask(id)
	if err != nil {
		return nil, err
	}
	if task.Status == constants.TaskStatusCompleted {
		return nil, utils.NewAppError(409, constants.TaskNotAcceptable,
			fmt.Sprintf(constants.TaskNotAcceptableMessage, id, task.Status))
	}
	now := time.Now().UTC()
	// 已签收过的任务保留首次签收时间（未签收任务才把 accepted_at 置为当前）。
	if task.AcceptedAt == nil {
		task.AcceptedAt = &now
	}
	task.Status = constants.TaskStatusAccepted
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	summary, _ := s.summaryOf(task.TurnaroundID)
	view := constructors.BuildGroundTaskView(*task, s.flightNo(task.TurnaroundID), summary, now)

	middlewares.QueueAudit(c, "TASK_ACCEPT", constants.TargetTask, fmt.Sprint(id),
		utils.RenderLog(constants.LogTemplates["TASK_ACCEPT"], map[string]interface{}{
			"TaskID":            id,
			"AcceptedAt":        now.Format(time.RFC3339),
			"EffectiveDeadline": view.EffectiveDeadline.Format(time.RFC3339),
		}))
	return &view, nil
}

// Complete 完成：actual_finish 默认取当前时间，超时按 actual_finish 比生效截止。
func (s *GroundTaskService) Complete(c *gin.Context, id uint) (*types.GroundTaskView, error) {
	task, err := s.loadTask(id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	if task.AcceptedAt == nil {
		task.AcceptedAt = &now
	}
	task.ActualFinish = &now
	task.Status = constants.TaskStatusCompleted
	task.BlockerNote = ""
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	summary, _ := s.summaryOf(task.TurnaroundID)
	view := constructors.BuildGroundTaskView(*task, s.flightNo(task.TurnaroundID), summary, now)
	overdueText := map[bool]string{true: fmt.Sprintf("是(%d 分钟)", view.OverdueMinutes), false: "否"}[view.Overdue]

	middlewares.QueueAudit(c, "TASK_COMPLETE", constants.TargetTask, fmt.Sprint(id),
		utils.RenderLog(constants.LogTemplates["TASK_COMPLETE"], map[string]interface{}{
			"TaskID":            id,
			"ActualFinish":      now.Format(time.RFC3339),
			"EffectiveDeadline": view.EffectiveDeadline.Format(time.RFC3339),
			"Overdue":           overdueText,
		}))
	return &view, nil
}

// Block 阻塞登记。
func (s *GroundTaskService) Block(c *gin.Context, id uint, req types.GroundTaskBlockRequest) (*types.GroundTaskView, error) {
	task, err := s.loadTask(id)
	if err != nil {
		return nil, err
	}
	if task.Status == constants.TaskStatusCompleted {
		return nil, utils.NewAppError(409, constants.TaskNotAcceptable,
			fmt.Sprintf(constants.TaskNotAcceptableMessage, id, task.Status))
	}
	task.Status = constants.TaskStatusBlocked
	task.BlockerNote = req.BlockerNote
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	summary, _ := s.summaryOf(task.TurnaroundID)
	view := constructors.BuildGroundTaskView(*task, s.flightNo(task.TurnaroundID), summary, time.Now().UTC())

	middlewares.QueueAudit(c, "TASK_BLOCK", constants.TargetTask, fmt.Sprint(id),
		utils.RenderLog(constants.LogTemplates["TASK_BLOCK"], map[string]interface{}{
			"TaskID": id, "BlockerNote": utils.OneLine(req.BlockerNote),
		}))
	return &view, nil
}
