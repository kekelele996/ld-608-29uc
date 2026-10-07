package services

import (
	"errors"
	"log"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
)

var (
	ErrTaskNotFound      = errors.New(constants.TaskNotFoundMessage)
	ErrTaskAlreadySigned = errors.New(constants.TaskAlreadySignedMessage)
)

func ListGroundTasks() []models.GroundTask {
	return repositories.ListGroundTasks()
}

// OpenDelayMinutes 统计航班未关闭延误的顺延分钟数。
// 口径与 frontend/src/utils/overtime.ts 保持一致：仅累加 resolved_at 为空（未关闭）的延误。
func OpenDelayMinutes(turnaroundID int, events []models.DelayEvent) int {
	total := 0
	for _, event := range events {
		if event.TurnaroundID == turnaroundID && event.ResolvedAt == "" {
			total += event.Minutes
		}
	}
	return total
}

// EffectiveDeadline 当前生效截止时间 = 任务原计划截止 + 未关闭延误分钟。
// 任务原 deadline 不改写，顺延结果仅另算返回；关闭延误后顺延自动失效。
func EffectiveDeadline(task models.GroundTask, events []models.DelayEvent) (time.Time, error) {
	base, err := time.Parse(time.RFC3339, task.Deadline)
	if err != nil {
		return time.Time{}, err
	}
	return base.Add(time.Duration(OpenDelayMinutes(task.TurnaroundID, events)) * time.Minute), nil
}

// SignOffGroundTask 签收未签收任务：签收时间写入 actual_finish，状态置为 SIGNED。
// 签收后任务超时判定改用 actual_finish 与生效截止比较。
func SignOffGroundTask(id int, signedAt string) (models.GroundTask, error) {
	task, ok := repositories.FindGroundTask(id)
	if !ok {
		return models.GroundTask{}, ErrTaskNotFound
	}
	if task.ActualFinish != "" {
		return models.GroundTask{}, ErrTaskAlreadySigned
	}
	if signedAt == "" {
		signedAt = time.Now().UTC().Format(time.RFC3339)
	}
	task.Status = constants.GroundTaskStatusSigned
	task.ActualFinish = signedAt
	repositories.UpdateGroundTask(task)
	log.Printf(constants.LogGroundTaskSignOff, task.ID, task.ActualFinish)
	return task, nil
}
