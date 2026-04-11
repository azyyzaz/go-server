package captcha

import (
	"go-server/internal/errcode"
	"go-server/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/captcha", h.Generate)
}

// Generate godoc
//
//	@Summary		获取图形验证码
//	@Description	生成登录验证码，返回验证码 ID 和图片内容。
//	@Tags			认证授权
//	@Produce		json
//	@Success		200		{object}	response.Body{data=GenerateResponse}
//	@Router			/auth/captcha [get]
func (h *Handler) Generate(c *gin.Context) {
	res, err := h.svc.Generate()
	if err != nil {
		_ = c.Error(errcode.ErrInternalError.AsError())
		return
	}
	response.Success(c, res)
}
