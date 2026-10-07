package middlewares

import (
	"strings"

	"groundTurn/src/constants"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// IdentityKey 认证身份在 gin.Context 中的键。
const IdentityKey = "identity"

// AuthMiddleware 校验 JWT 并注入身份。/api/auth/login 与 /health 不挂该中间件。
func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer"))
		token = strings.TrimSpace(token)
		if token == "" {
			utils.Fail(c, 401, constants.AuthRequired, constants.AuthRequiredMessage)
			c.Abort()
			return
		}
		id, err := utils.ParseJWT(token, secret)
		if err != nil {
			utils.Fail(c, 401, constants.AuthRequired, constants.AuthRequiredMessage)
			c.Abort()
			return
		}
		c.Set(IdentityKey, id)
		c.Next()
	}
}

// CurrentIdentity 从上下文取身份，路由处理器复用。
func CurrentIdentity(c *gin.Context) *utils.AuthIdentity {
	if raw, ok := c.Get(IdentityKey); ok {
		return raw.(*utils.AuthIdentity)
	}
	return nil
}
