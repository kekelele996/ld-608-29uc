package constructors

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/types"
)

func NewResourceBooking(req types.ResourceBookingRequest) models.ResourceBooking {
	return models.ResourceBooking{
		ResourceID:    req.ResourceID,
		TurnaroundID:  req.TurnaroundID,
		TaskID:        req.TaskID,
		StartTime:     req.StartTime.UTC(),
		EndTime:       req.EndTime.UTC(),
		BookingStatus: constants.BookingHeld,
	}
}

func BuildResourceBookingView(b models.ResourceBooking, resourceCode string) types.ResourceBookingView {
	return types.ResourceBookingView{
		ID:             b.ID,
		ResourceID:     b.ResourceID,
		ResourceCode:   resourceCode,
		TurnaroundID:   b.TurnaroundID,
		TaskID:         b.TaskID,
		StartTime:      b.StartTime,
		EndTime:        b.EndTime,
		BookingStatus:  b.BookingStatus,
		StatusText:     constants.StatusText["BookingStatus"][b.BookingStatus],
		ConflictReason: b.ConflictReason,
	}
}
