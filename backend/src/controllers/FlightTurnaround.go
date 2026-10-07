package controllers

import (
	"strconv"

	"groundTurn/src/services"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type FlightTurnaroundController struct {
	service *services.FlightTurnaroundService
}

func NewFlightTurnaroundController(service *services.FlightTurnaroundService) *FlightTurnaroundController {
	return &FlightTurnaroundController{service: service}
}

func (ctl *FlightTurnaroundController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/flight-turnarounds", ctl.Register)
	rg.POST("/flight-turnarounds/:id/status", ctl.UpdateStatus)
}

func (ctl *FlightTurnaroundController) Register(c *gin.Context) {
	var req types.FlightTurnaroundRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.Register(c, req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Created(c, view)
}

func (ctl *FlightTurnaroundController) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", "invalid turnaround id")
		return
	}
	var req types.FlightTurnaroundStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.UpdateStatus(c, uint(id), req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}
