package controllers

import (
	"groundTurn/src/services"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type GroundResourceController struct {
	service *services.GroundResourceService
}

func NewGroundResourceController(service *services.GroundResourceService) *GroundResourceController {
	return &GroundResourceController{service: service}
}

func (ctl *GroundResourceController) RegisterReadRoutes(rg *gin.RouterGroup) {
	rg.GET("/ground-resources", ctl.List)
}

func (ctl *GroundResourceController) RegisterWriteRoutes(rg *gin.RouterGroup) {
	rg.POST("/ground-resources", ctl.Upsert)
}

func (ctl *GroundResourceController) List(c *gin.Context) {
	views, err := ctl.service.List()
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, views)
}

func (ctl *GroundResourceController) Upsert(c *gin.Context) {
	var req types.GroundResourceUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.Upsert(req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.Created(c, view)
}
