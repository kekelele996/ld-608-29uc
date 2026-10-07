package controllers

import (
	"net/http"

	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// handleError controller 层二次包装业务错误，保留错误码（不在这里吞错）。
func handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*utils.AppError); ok {
		utils.Fail(c, appErr.Status, appErr.Code, appErr.Message)
		return
	}
	utils.Fail(c, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}
