package http

import (
	"github.com/mojocn/base64Captcha"
	"net/http"

	"backend/pkg/errors"
	"backend/pkg/response"
	"github.com/gin-gonic/gin"
)

var store = base64Captcha.DefaultMemStore

// CaptchaHandler 处理验证码相关
type CaptchaHandler struct{}

func NewCaptchaHandler(r *gin.Engine) {
	handler := &CaptchaHandler{}
	r.GET("/api/v1/auth/captcha", handler.GenerateCaptcha)
}

func (h *CaptchaHandler) GenerateCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusOK, response.Error(errors.ErrInternalServer))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"captcha_id": id,
		"image_url":  b64s,
	}))
}

// VerifyCaptcha 验证验证码
func VerifyCaptcha(id string, answer string) bool {
	return store.Verify(id, answer, true)
}
