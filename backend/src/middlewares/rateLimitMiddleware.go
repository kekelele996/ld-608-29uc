package middlewares

import (
	"sync"
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/utils"

	"github.com/gin-gonic/gin"
)

// 简易单 IP 令牌/固定窗口限流，本地演示足够。
type rateBucket struct {
	windowStart time.Time
	count       int
}

// RateLimitMiddleware 每个 IP 每 window 最多 limit 次。
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]*rateBucket{}
	return func(c *gin.Context) {
		key := c.ClientIP()
		mu.Lock()
		b, ok := buckets[key]
		now := time.Now()
		if !ok || now.Sub(b.windowStart) > window {
			b = &rateBucket{windowStart: now}
			buckets[key] = b
		}
		b.count++
		allowed := b.count <= limit
		mu.Unlock()
		if !allowed {
			utils.Fail(c, 429, constants.RateLimited, constants.RateLimitedMessage)
			c.Abort()
			return
		}
		c.Next()
	}
}
