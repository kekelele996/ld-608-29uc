package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"groundTurn/src/constants"
	"groundTurn/src/services"
	"groundTurn/src/types"
)

func ListDelayEvent(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListDelayEvents())
}

func CreateDelayEvent(c *gin.Context) {
	var req types.CreateDelayEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constants.ValidationFailed, "message": constants.ValidationFailedMessage})
		return
	}
	event, err := services.RegisterDelayEvent(req)
	if err != nil {
		if errors.Is(err, services.ErrDelayEventInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"code": constants.ValidationFailed, "message": constants.ValidationFailedMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func CloseDelayEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": constants.ValidationFailed, "message": constants.ValidationFailedMessage})
		return
	}
	var req types.CloseDelayEventRequest
	_ = c.ShouldBindJSON(&req) // 允许空 body，closed_at 缺省由服务端取当前时间
	event, err := services.CloseDelayEvent(id, req.ClosedAt)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrDelayEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"code": constants.DelayEventNotFound, "message": constants.DelayEventNotFoundMessage})
		case errors.Is(err, services.ErrDelayEventAlreadyClosed):
			c.JSON(http.StatusConflict, gin.H{"code": constants.DelayEventAlreadyClosed, "message": constants.DelayEventAlreadyClosedMessage})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"code": "INTERNAL", "message": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, event)
}
