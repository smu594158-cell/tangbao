package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"backend/pkg/errors"
	"backend/pkg/response"
)

func realIP(c *gin.Context) string {
	if fwd := c.GetHeader("X-Forwarded-For"); fwd != "" {
		if idx := strings.Index(fwd, ","); idx > 0 {
			return strings.TrimSpace(fwd[:idx])
		}
		return strings.TrimSpace(fwd)
	}
	if real := c.GetHeader("X-Real-IP"); real != "" {
		return real
	}
	return c.ClientIP()
}

type visitorEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IP rate limiter map with cleanup
var (
	visitors   = make(map[string]*visitorEntry)
	mu         sync.Mutex
	cleanerOnce sync.Once
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()

	entry, exists := visitors[ip]
	if !exists {
		entry = &visitorEntry{
			limiter:  rate.NewLimiter(5, 10),
			lastSeen: time.Now(),
		}
		visitors[ip] = entry
	} else {
		entry.lastSeen = time.Now()
	}
	mu.Unlock()

	// 启动定期清理（仅一次）
	cleanerOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				mu.Lock()
				now := time.Now()
				for ip, entry := range visitors {
					if now.Sub(entry.lastSeen) > 10*time.Minute {
						delete(visitors, ip)
					}
				}
				mu.Unlock()
			}
		}()
	})

	return entry.limiter
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getVisitor(realIP(c))
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, response.Error(errors.New(42900, "操作太频繁，请稍后再试")))
			c.Abort()
			return
		}
		c.Next()
	}
}
