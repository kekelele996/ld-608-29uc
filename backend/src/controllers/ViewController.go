package controllers

import (
	"groundTurn/src/services"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// ViewController 暴露读模型列表。任务列表与看板共用 ViewService 的同一套顺延/超时口径。
type ViewController struct{ viewService *services.ViewService }

func NewViewController(vs *services.ViewService) *ViewController {
	return &ViewController{viewService: vs}
}

func (ctl *ViewController) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ground-tasks", ctl.ListTasks)
	rg.GET("/flight-turnarounds", ctl.ListTurnarounds)
}

func (ctl *ViewController) ListTasks(c *gin.Context) {
	bundle, err := ctl.viewService.LoadBundle()
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, bundle.Tasks)
}

func (ctl *ViewController) ListTurnarounds(c *gin.Context) {
	bundle, err := ctl.viewService.LoadBundle()
	if err != nil {
		handleError(c, err)
		return
	}
	utils.OK(c, bundle.Turnarounds)
}
