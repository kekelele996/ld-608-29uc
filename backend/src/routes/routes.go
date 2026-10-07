package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"groundTurn/src/controllers"
)

func Start(addr string) {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ground-turn"}) })
	r.GET("/api/flight-turnaround", controllers.ListFlightTurnaround)
	r.GET("/api/ground-task", controllers.ListGroundTask)
	r.POST("/api/ground-task/:id/sign-off", controllers.SignOffGroundTask)
	r.GET("/api/delay-event", controllers.ListDelayEvent)
	r.POST("/api/delay-event", controllers.CreateDelayEvent)
	r.POST("/api/delay-event/:id/close", controllers.CloseDelayEvent)
	r.GET("/api/ground-resource", func(c *gin.Context) { c.JSON(http.StatusOK, []gin.H{{"id": 1, "name": "保障资源", "status": "READY"}}) })
	r.GET("/api/resource-booking", func(c *gin.Context) { c.JSON(http.StatusOK, []gin.H{{"id": 1, "name": "资源预约", "status": "READY"}}) })
	r.Run(addr)
}
