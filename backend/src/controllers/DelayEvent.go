package controllers

import (
	"strconv"

	"groundTurn/src/services"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type DelayEventController struct{ service *services.DelayEventService }

func NewDelayEventController(service *services.DelayEventService) *DelayEventController {
	return &DelayEventController{service: service}
}

func (ctl *DelayEventController) RegisterReadRoutes(rg *gin.RouterGroup) {
	rg.GET("/delay-events", ctl.List)
}

func (ctl *DelayEventController) RegisterWriteRoutes(rg *gin.RouterGroup) {
	rg.POST("/delay-events", ctl.Register)
	rg.POST("/delay-events/:id/resolve", ctl.Resolve)
}

func (ctl *DelayEventController) List(c *gin.Context) {
	views, err := ctl.service.List()
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, views)
}

func (ctl *DelayEventController) Register(c *gin.Context) {
	var req types.DelayEventRegisterRequest
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

func (ctl *DelayEventController) Resolve(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", "invalid delay id")
		return
	}
	view, err := ctl.service.Resolve(c, uint(id))
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}
