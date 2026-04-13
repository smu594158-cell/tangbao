package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"backend/pkg/errors"
	"backend/pkg/response"
)

// IP rate limiter map
var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// 允许每秒 5 个请求，突发最多 10 个
		limiter = rate.NewLimiter(5, 10)
		visitors[ip] = limiter
	}

	return limiter
}

// RateLimit 限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getVisitor(c.ClientIP())
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, response.Error(errors.New(42900, "操作太频繁，请稍后再试")))
			c.Abort()
			return
		}
		c.Next()
	}
}
