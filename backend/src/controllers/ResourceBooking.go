package controllers

import (
	"strconv"

	"groundTurn/src/services"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type ResourceBookingController struct {
	service *services.ResourceBookingService
}

func NewResourceBookingController(service *services.ResourceBookingService) *ResourceBookingController {
	return &ResourceBookingController{service: service}
}

func (ctl *ResourceBookingController) RegisterReadRoutes(rg *gin.RouterGroup) {
	rg.GET("/resource-bookings", ctl.List)
}

func (ctl *ResourceBookingController) RegisterWriteRoutes(rg *gin.RouterGroup) {
	rg.POST("/resource-bookings", ctl.Create)
	rg.POST("/resource-bookings/:id/release", ctl.Release)
}

func (ctl *ResourceBookingController) List(c *gin.Context) {
	views, err := ctl.service.List()
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, views)
}

func (ctl *ResourceBookingController) Create(c *gin.Context) {
	var req types.ResourceBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.Create(c, req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Created(c, view)
}

func (ctl *ResourceBookingController) Release(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", "invalid booking id")
		return
	}
	view, err := ctl.service.Release(c, uint(id))
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}
