package controllers

import (
	"strconv"

	"groundTurn/src/services"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct{ service *services.AuthService }

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (ctl *AuthController) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/login", ctl.Login)
}

func (ctl *AuthController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/audit-logs", ctl.ListLogs)
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, 400, "VALIDATION_FAILED", err.Error())
		return
	}
	view, err := ctl.service.Login(c, req)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}

func (ctl *AuthController) ListLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	views, err := ctl.service.ListLogs(limit)
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, views)
}

// DashboardController 看板汇总。
type DashboardController struct{ service *services.DashboardService }

func NewDashboardController(service *services.DashboardService) *DashboardController {
	return &DashboardController{service: service}
}

func (ctl *DashboardController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/dashboard", ctl.Dashboard)
}

func (ctl *DashboardController) Dashboard(c *gin.Context) {
	view, err := ctl.service.Build()
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, view)
}
