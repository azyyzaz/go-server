package profile

import (
	"go-server/internal/errcode"
	"go-server/internal/middleware"
	"go-server/internal/modules/audit"
	"go-server/internal/response"
	"go-server/internal/validation"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.GetProfile)
	rg.PUT("", h.UpdateProfile)
	rg.PUT("/password", h.ChangePassword)
	rg.POST("/avatar", h.UploadAvatar)
	rg.GET("/login-logs", h.ListMyLoginLogs)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	result, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}

	result, err := h.svc.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) ChangePassword(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}

	if err := h.svc.ChangePassword(c.Request.Context(), userID, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"changed": true})
}

func (h *Handler) UploadAvatar(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}

	result, err := h.svc.UploadAvatar(c.Request.Context(), userID, file)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) ListMyLoginLogs(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}

	var q audit.ListLoginLogsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	result, err := h.svc.ListMyLoginLogs(c.Request.Context(), userID, q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func currentUserID(c *gin.Context) (uint, bool) {
	raw, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		_ = c.Error(errcode.ErrUnauthorized.AsError())
		return 0, false
	}

	userID, ok := raw.(uint)
	if !ok {
		_ = c.Error(errcode.ErrUnauthorized.AsError())
		return 0, false
	}

	return userID, true
}
