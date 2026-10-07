package constructors

import (
	"groundTurn/src/models"
	"groundTurn/src/types"
)

// NewDelayEventFromRequest 由登记请求构造未关闭的延误事件（resolved_at 为空即参与顺延）。
func NewDelayEventFromRequest(id int, req types.CreateDelayEventRequest) models.DelayEvent {
	return models.DelayEvent{
		ID:                 id,
		TurnaroundID:       req.TurnaroundID,
		DelayType:          req.DelayType,
		Minutes:            req.Minutes,
		RootCause:          req.RootCause,
		ResponsibilityTeam: req.ResponsibilityTeam,
		ResolvedAt:         "",
	}
}
