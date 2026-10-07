package middlewares

import (
	"time"

	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// 写操作经 Audit 记录一条日志。业务模板渲染在 service 内完成并放进上下文，
// 中间件负责统一落库（跨切面：控制器不直接碰 audit_log 表）。
const auditDetailKey = "audit_detail"

type auditPayload struct {
	Action     string
	TargetType string
	TargetID   string
	Detail     string
}

// QueueAudit service/controller 登记待写日志。
func QueueAudit(c *gin.Context, action, targetType, targetID, detail string) {
	c.Set(auditDetailKey, auditPayload{Action: action, TargetType: targetType, TargetID: targetID, Detail: detail})
}

// AuditLogMiddleware 在请求成功返回后把排队的日志写入 audit_log。
func AuditLogMiddleware(repo *repositories.AuditLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 400 {
			return
		}
		raw, ok := c.Get(auditDetailKey)
		if !ok {
			return
		}
		payload := raw.(auditPayload)
		actor := "anonymous"
		if id, exists := c.Get(IdentityKey); exists {
			actor = id.(*utils.AuthIdentity).Username
		}
		_ = repo.Create(&models.AuditLog{
			Actor:      actor,
			Action:     payload.Action,
			TargetType: payload.TargetType,
			TargetID:   payload.TargetID,
			Detail:     payload.Detail,
			CreatedAt:  time.Now().UTC(),
		})
	}
}
