package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// OK / Fail 统一响应封装，controller 只通过这里出参。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "data": data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"ok": true, "data": data})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, gin.H{"ok": false, "error": gin.H{"code": code, "message": message}})
}
