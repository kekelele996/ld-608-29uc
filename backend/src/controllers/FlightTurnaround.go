package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/services"
)

func ListFlightTurnaround(c *gin.Context) {
	c.JSON(http.StatusOK, services.ListFlightTurnarounds())
}
