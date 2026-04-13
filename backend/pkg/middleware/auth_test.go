package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"backend/pkg/middleware"
	"backend/pkg/utils"
)

func TestAdminAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	secret := "test_secret"
	
	// 设置测试路由
	r := gin.New()
	r.Use(middleware.JWTAuth(secret))
	r.Use(middleware.AdminAuth())
	r.GET("/admin/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 辅助函数：生成 Token
	generateTestToken := func(role int8) string {
		token, _ := utils.GenerateToken(1, "testuser", role, secret, time.Hour)
		return token
	}

	tests := []struct {
		name         string
		role         int8
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Normal User Access",
			role:         1, // domain.RoleUser
			expectedCode: http.StatusOK, // 业务错误码包装在200中
			expectedBody: `{"code":40007,"message":"禁止访问：权限不足"}`,
		},
		{
			name:         "Admin User Access",
			role:         9, // domain.RoleAdmin
			expectedCode: http.StatusOK,
			expectedBody: `{"message":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := generateTestToken(tt.role)
			req, _ := http.NewRequest("GET", "/admin/test", nil)
			req.Header.Set("Authorization", "Bearer "+token)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("expected HTTP status %d, got %d", tt.expectedCode, w.Code)
			}
			
			if w.Body.String() != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, w.Body.String())
			}
		})
	}
}
