package constructors

import (
	"groundTurn/src/models"
)

// GroundTaskResponse 在实体字段之外附带顺延口径的计算结果：
// deadline 保持原计划时间不改写，顺延结果以 effective_deadline 另算返回。
type GroundTaskResponse struct {
	models.GroundTask
	OpenDelayMinutes  int    `json:"open_delay_minutes"`
	EffectiveDeadline string `json:"effective_deadline"`
}

func NewGroundTaskResponse(task models.GroundTask, openDelayMinutes int, effectiveDeadline string) GroundTaskResponse {
	return GroundTaskResponse{
		GroundTask:        task,
		OpenDelayMinutes:  openDelayMinutes,
		EffectiveDeadline: effectiveDeadline,
	}
}
