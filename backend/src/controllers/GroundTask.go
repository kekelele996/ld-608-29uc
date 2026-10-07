package controllers

import (
	"strconv"

	"groundTurn/src/services"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type GroundTaskController struct{ service *services.GroundTaskService }

func NewGroundTaskController(service *services.GroundTaskService) *GroundTaskController {
	return &GroundTaskController{service: service}
}

// RegisterDispatchRoutes 派工（调度/督导）。
func (ctl *GroundTaskController) RegisterDispatchRoutes(rg *gin.RouterGroup) {
	rg.POST("/ground-tasks", ctl.Dispatch)
}

// RegisterActionRoutes 签收/完成/阻塞（班组/调度/督导）。
func (ctl *GroundTaskController) RegisterActionRoutes(rg *gin.RouterGroup) {
	rg.POST("/ground-tasks/:id/accept", ctl.Accept)
	rg.POST("/ground-tasks/:id/complete", ctl.Complete)
	rg.POST("/ground-tasks/:id/block", ctl.Block)
}

func (ctl *GroundTaskController) Dispatch(c *gin.Context) {
	var req types.GroundTaskDispatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.Dispatch(c, req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Created(c, view)
}

func (ctl *GroundTaskController) Accept(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", "invalid task id")
		return
	}
	view, err := ctl.service.Accept(c, uint(id))
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}

func (ctl *GroundTaskController) Complete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", "invalid task id")
		return
	}
	view, err := ctl.service.Complete(c, uint(id))
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}

func (ctl *GroundTaskController) Block(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", "invalid task id")
		return
	}
	var req types.GroundTaskBlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.Block(c, uint(id), req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}
