package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/pkg/errors"
	"backend/pkg/response"
	"backend/pkg/utils"
)

// JWTAuth 中间件，验证携带的token
func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, response.Error(errors.ErrUnauthorized))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, response.Error(errors.ErrTokenInvalid))
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], secret)
		if err != nil {
			c.JSON(http.StatusOK, response.Error(errors.ErrTokenInvalid))
			c.Abort()
			return
		}

		// 将当前请求的 userID 信息保存到请求的上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next() // 后续的处理函数可以用c.Get("userID") 来获取当前请求的用户信息
	}
}

// AdminAuth 管理员权限校验中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusOK, response.Error(errors.ErrUnauthorized))
			c.Abort()
			return
		}

		role, ok := roleVal.(int8)
		if !ok || role != 9 { // domain.RoleAdmin = 9
			c.JSON(http.StatusOK, response.Error(errors.ErrForbidden))
			c.Abort()
			return
		}

		c.Next()
	}
}
