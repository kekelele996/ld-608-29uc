package middlewares

import (
	"net/http"

	"groundTurn/src/constants"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware 统一兜底：AppError 按其状态码返回，其余记为内部错误。
// service/controller 已分别包装过业务错误，这里只做最后一道，不吞错误码。
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		if appErr, ok := err.(*utils.AppError); ok {
			utils.Fail(c, appErr.Status, appErr.Code, appErr.Message)
			return
		}
		utils.Fail(c, http.StatusInternalServerError, constants.InternalError, constants.InternalErrorMessage)
	}
}
