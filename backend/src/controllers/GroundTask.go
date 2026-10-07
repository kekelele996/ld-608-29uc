package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/constructors"
	"groundTurn/src/models"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

// groundTaskResponse 组装带顺延口径的响应：原 deadline 不改写，effective_deadline 另算。
func groundTaskResponse(task models.GroundTask, events []models.DelayEvent) constructors.GroundTaskResponse {
	effectiveText := task.Deadline
	if effective, err := services.EffectiveDeadline(task, events); err == nil {
		effectiveText = effective.UTC().Format(time.RFC3339)
	}
	return constructors.NewGroundTaskResponse(task, services.OpenDelayMinutes(task.TurnaroundID, events), effectiveText)
}

func ListGroundTask(c *gin.Context) {
	tasks := services.ListGroundTasks()
	events := services.ListDelayEvents()
	resp := make([]constructors.GroundTaskResponse, 0, len(tasks))
	for _, task := range tasks {
		resp = append(resp, groundTaskResponse(task, events))
	}
	c.JSON(http.StatusOK, resp)
}

func SignOffGroundTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constants.ValidationFailed, "message": constants.ValidationFailedMessage})
		return
	}
	var req types.SignOffGroundTaskRequest
	_ = c.ShouldBindJSON(&req) // 允许空 body，signed_at 缺省由服务端取当前时间
	task, err := services.SignOffGroundTask(id, req.SignedAt)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": constants.TaskNotFound, "message": constants.TaskNotFoundMessage})
		case errors.Is(err, services.ErrTaskAlreadySigned):
			c.JSON(http.StatusConflict, gin.H{"code": constants.TaskAlreadySigned, "message": constants.TaskAlreadySignedMessage})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, groundTaskResponse(task, services.ListDelayEvents()))
}
