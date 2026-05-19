package middleware

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"

	"backend/internal/domain"
	"backend/internal/repository"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AdminOperationLog 记录管理员操作日志
func AdminOperationLog(repo repository.AdminLogRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅记录修改操作 (POST, PUT, DELETE)
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		// 获取用户信息
		userIDVal, exists := c.Get("userID")
		usernameVal, _ := c.Get("username")
		if !exists {
			c.Next()
			return
		}
		userID := userIDVal.(uint64)
		username := usernameVal.(string)

		action := "unknown"
		if strings.Contains(c.Request.URL.Path, "/admin/users") {
			action = "manage_user"
		} else if strings.Contains(c.Request.URL.Path, "/admin/content") {
			action = "manage_content"
		}

		logEntry := &domain.AdminLog{
			AdminID:  userID,
			Username: username,
			Action:   action,
			Method:   c.Request.Method,
			Path:     c.Request.URL.Path,
			IP:       c.ClientIP(),
		}

		// 异步记录日志
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[AdminLog] recovered from panic: %v", r)
				}
			}()
			if err := repo.Create(context.Background(), logEntry); err != nil {
				log.Printf("[AdminLog] failed to save admin log: %v", err)
			}
		}()

		c.Next()
	}
}
