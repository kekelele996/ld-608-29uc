package constructors

import (
	"time"

	"groundTurn/src/models"
	"groundTurn/src/types"
	"groundTurn/src/utils"
)

// NewDelayEvent 登记延误：创建时未关闭（resolved_at=nil），立即进入顺延口径。
func NewDelayEvent(req types.DelayEventRegisterRequest, now time.Time) models.DelayEvent {
	return models.DelayEvent{
		TurnaroundID:       req.TurnaroundID,
		DelayType:          req.DelayType,
		Minutes:            req.Minutes,
		RootCause:          req.RootCause,
		ResponsibilityTeam: req.ResponsibilityTeam,
		CreatedAt:          now,
		ResolvedAt:         nil,
	}
}

// BuildDelayEventView 组装延误响应，并带出该航班当前未关闭延误累计分钟。
func BuildDelayEventView(e models.DelayEvent, flightNo string, summary utils.OpenDelaySummary) types.DelayEventView {
	return types.DelayEventView{
		ID:                 e.ID,
		TurnaroundID:       e.TurnaroundID,
		FlightNo:           flightNo,
		DelayType:          e.DelayType,
		Minutes:            e.Minutes,
		RootCause:          e.RootCause,
		ResponsibilityTeam: e.ResponsibilityTeam,
		CreatedAt:          e.CreatedAt,
		ResolvedAt:         e.ResolvedAt,
		Closed:             e.ResolvedAt != nil,
		OpenDelayMinutes:   summary.OpenMinutes,
	}
}
