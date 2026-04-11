package auth

import (
	"go-server/internal/errcode"
	"go-server/internal/modules/audit"
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
	rg.POST("/login", h.Login)
	rg.POST("/logout", h.Logout)
	rg.POST("/refresh", h.Refresh)
}

// Login godoc
//
//	@Summary		用户登录
//	@Description	使用用户名、密码和验证码完成登录，成功后返回 access_token 和 refresh_token。
//	@Tags			认证授权
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"登录请求参数"
//	@Success		200		{object}	response.Body{data=TokenResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Router			/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	ctx := audit.WithLoginMeta(c.Request.Context(), audit.LoginMeta{
		RequestID: c.GetString("request_id"),
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	})

	tokens, err := h.svc.Login(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, tokens)
}

// Logout godoc
//
//	@Summary		用户登出
//	@Description	将当前 access_token 加入黑名单，使其立即失效。
//	@Tags			认证授权
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LogoutRequest	true	"登出请求参数"
//	@Success		200		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	if err := h.svc.Logout(c.Request.Context(), req.AccessToken); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, nil)
}

// Refresh godoc
//
//	@Summary		刷新令牌
//	@Description	使用 refresh_token 换取一组新的 access_token 和 refresh_token。
//	@Tags			认证授权
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshRequest	true	"刷新令牌请求参数"
//	@Success		200		{object}	response.Body{data=TokenResponse}
//	@Failure		401		{object}	response.Body
//	@Router			/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	tokens, err := h.svc.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, tokens)
}
