package services

import (
	"errors"
	"log"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
)

var (
	ErrDelayEventNotFound      = errors.New(constants.DelayEventNotFoundMessage)
	ErrDelayEventAlreadyClosed = errors.New(constants.DelayEventAlreadyClosedMessage)
	ErrDelayEventInvalid       = errors.New(constants.ValidationFailedMessage)
)

func ListDelayEvents() []models.DelayEvent {
	return repositories.ListDelayEvents()
}

// RegisterDelayEvent 登记延误：新事件 resolved_at 为空（未关闭），
// 登记后立即参与该航班任务的截止时间顺延（见 EffectiveDeadline）。
func RegisterDelayEvent(req types.CreateDelayEventRequest) (models.DelayEvent, error) {
	if req.TurnaroundID <= 0 || req.Minutes <= 0 {
		return models.DelayEvent{}, ErrDelayEventInvalid
	}
	event := constructors.NewDelayEventFromRequest(repositories.NextDelayEventID(), req)
	repositories.CreateDelayEvent(event)
	log.Printf(constants.LogDelayEventRegister, event.ID, event.TurnaroundID, event.Minutes)
	return event, nil
}

// CloseDelayEvent 关闭延误：写入 resolved_at，关闭后不再参与任务截止顺延。
func CloseDelayEvent(id int, closedAt string) (models.DelayEvent, error) {
	event, ok := repositories.FindDelayEvent(id)
	if !ok {
		return models.DelayEvent{}, ErrDelayEventNotFound
	}
	if event.ResolvedAt != "" {
		return models.DelayEvent{}, ErrDelayEventAlreadyClosed
	}
	if closedAt == "" {
		closedAt = time.Now().UTC().Format(time.RFC3339)
	}
	event.ResolvedAt = closedAt
	repositories.UpdateDelayEvent(event)
	log.Printf(constants.LogDelayEventClose, event.ID, event.ResolvedAt)
	return event, nil
}
