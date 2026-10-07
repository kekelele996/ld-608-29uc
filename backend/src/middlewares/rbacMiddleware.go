package middlewares

import (
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware 校验角色是否在允许名单。触达：登录角色、路由守卫、按钮显隐（前端）。
func RBACMiddleware(allowed ...string) gin.HandlerFunc {
	allow := map[string]bool{}
	for _, role := range allowed {
		allow[role] = true
	}
	return func(c *gin.Context) {
		id := CurrentIdentity(c)
		if id == nil {
			utils.Fail(c, 401, "AUTH_REQUIRED", "missing token")
			c.Abort()
			return
		}
		if !allow[id.Role] {
			action := c.Request.Method + " " + c.FullPath()
			appErr := utils.ForbiddenError(id.Role, action)
			utils.Fail(c, appErr.Status, appErr.Code, appErr.Message)
			c.Abort()
			return
		}
		c.Next()
	}
}
