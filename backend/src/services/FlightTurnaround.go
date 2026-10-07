package services

import (
	"fmt"
	"slices"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/middlewares"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type FlightTurnaroundService struct {
	repo        *repositories.FlightTurnaroundRepository
	viewService *ViewService
}

func NewFlightTurnaroundService(repo *repositories.FlightTurnaroundRepository, vs *ViewService) *FlightTurnaroundService {
	return &FlightTurnaroundService{repo: repo, viewService: vs}
}

func (s *FlightTurnaroundService) Register(c *gin.Context, req types.FlightTurnaroundRegisterRequest) (*types.FlightTurnaroundView, error) {
	if req.DepartureTime.Before(req.ArrivalTime) {
		return nil, utils.ValidationError("departure_time must be after arrival_time")
	}
	flight := constructors.NewFlightTurnaround(req)
	if err := s.repo.Create(&flight); err != nil {
		return nil, err
	}
	view := constructors.BuildFlightTurnaroundView(flight, nil, utils.OpenDelaySummary{})
	middlewares.QueueAudit(c, "TURNAROUND_REGISTER", constants.TargetTurnaround, fmt.Sprint(flight.ID),
		utils.RenderLog(constants.LogTemplates["TURNAROUND_REGISTER"], map[string]interface{}{
			"FlightNo": flight.FlightNo, "StandNo": flight.StandNo,
			"DepartureTime": flight.DepartureTime.Format(time.RFC3339),
		}))
	return &view, nil
}

func (s *FlightTurnaroundService) UpdateStatus(c *gin.Context, id uint, req types.FlightTurnaroundStatusRequest) (*types.FlightTurnaroundView, error) {
	if !slices.Contains(constants.TurnaroundStatus, req.Status) {
		return nil, utils.EnumError("turnaround_status", req.Status)
	}
	flight, err := s.repo.Get(id)
	if err != nil {
		return nil, utils.NotFoundError(constants.TargetTurnaround, id)
	}
	from := flight.TurnaroundStatus
	if err := s.repo.UpdateStatus(id, req.Status, req.DelayReason); err != nil {
		return nil, err
	}
	bundle, err := s.viewService.LoadBundle()
	if err != nil {
		return nil, err
	}
	var view types.FlightTurnaroundView
	for _, tv := range bundle.Turnarounds {
		if tv.ID == id {
			view = tv
		}
	}
	middlewares.QueueAudit(c, "TURNAROUND_UPDATE", constants.TargetTurnaround, fmt.Sprint(id),
		utils.RenderLog(constants.LogTemplates["TURNAROUND_UPDATE"], map[string]interface{}{
			"FlightNo": flight.FlightNo, "FromStatus": from, "ToStatus": req.Status,
		}))
	return &view, nil
}
