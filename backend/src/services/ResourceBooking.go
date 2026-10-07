package services

import (
	"fmt"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/middlewares"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type ResourceBookingService struct {
	repo         *repositories.ResourceBookingRepository
	resourceRepo *repositories.GroundResourceRepository
}

func NewResourceBookingService(repo *repositories.ResourceBookingRepository, rr *repositories.GroundResourceRepository) *ResourceBookingService {
	return &ResourceBookingService{repo: repo, resourceRepo: rr}
}

func (s *ResourceBookingService) List() ([]types.ResourceBookingView, error) {
	rows, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	views := make([]types.ResourceBookingView, 0, len(rows))
	for i := range rows {
		code := ""
		if r, ferr := s.resourceRepo.Get(rows[i].ResourceID); ferr == nil {
			code = r.ResourceCode
		}
		views = append(views, constructors.BuildResourceBookingView(rows[i], code))
	}
	return views, nil
}

// Create 预约时做同资源时间窗冲突检测，冲突预约以 CONFLICT 状态落库并给出原因。
func (s *ResourceBookingService) Create(c *gin.Context, req types.ResourceBookingRequest) (*types.ResourceBookingView, error) {
	if !req.EndTime.After(req.StartTime) {
		return nil, utils.ValidationError("end_time must be after start_time")
	}
	resource, err := s.resourceRepo.Get(req.ResourceID)
	if err != nil {
		return nil, utils.NotFoundError(constants.TargetResource, req.ResourceID)
	}
	booking := constructors.NewResourceBooking(req)
	overlaps, err := s.repo.FindOverlap(req.ResourceID, req.StartTime, req.EndTime, 0)
	if err != nil {
		return nil, err
	}
	if len(overlaps) > 0 {
		booking.BookingStatus = constants.BookingConflict
		booking.ConflictReason = fmt.Sprintf("overlaps booking #%d", overlaps[0].ID)
	}
	if err := s.repo.Create(&booking); err != nil {
		return nil, err
	}
	view := constructors.BuildResourceBookingView(booking, resource.ResourceCode)

	if len(overlaps) > 0 {
		middlewares.QueueAudit(c, "BOOKING_CONFLICT", constants.TargetBooking, fmt.Sprint(booking.ID),
			utils.RenderLog(constants.LogTemplates["BOOKING_CONFLICT"], map[string]interface{}{
				"ResourceID": req.ResourceID, "ConflictID": overlaps[0].ID,
				"ConflictReason": booking.ConflictReason,
			}))
	} else {
		middlewares.QueueAudit(c, "BOOKING_CREATE", constants.TargetBooking, fmt.Sprint(booking.ID),
			utils.RenderLog(constants.LogTemplates["BOOKING_CREATE"], map[string]interface{}{
				"BookingID": booking.ID, "ResourceID": req.ResourceID,
				"StartTime": req.StartTime.Format("2006-01-02T15:04:05Z07:00"),
				"EndTime":   req.EndTime.Format("2006-01-02T15:04:05Z07:00"),
			}))
	}
	return &view, nil
}

// Release 释放预约。
func (s *ResourceBookingService) Release(c *gin.Context, id uint) (*types.ResourceBookingView, error) {
	booking, err := s.repo.Get(id)
	if err != nil {
		return nil, utils.NotFoundError(constants.TargetBooking, id)
	}
	if err := s.repo.UpdateStatus(id, constants.BookingReleased, booking.ConflictReason); err != nil {
		return nil, err
	}
	booking.BookingStatus = constants.BookingReleased
	code := ""
	if r, ferr := s.resourceRepo.Get(booking.ResourceID); ferr == nil {
		code = r.ResourceCode
	}
	view := constructors.BuildResourceBookingView(*booking, code)
	middlewares.QueueAudit(c, "BOOKING_RELEASE", constants.TargetBooking, fmt.Sprint(id),
		utils.RenderLog(constants.LogTemplates["BOOKING_RELEASE"], map[string]interface{}{"BookingID": id}))
	return &view, nil
}
