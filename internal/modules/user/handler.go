package user

import (
	"strconv"

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
	rg.GET("", h.ListUsers)
	rg.GET("/export", h.ExportUsers)
	rg.POST("", h.CreateUser)
	rg.POST("/batch-delete", h.BatchDeleteUsers)
	rg.GET("/:id", h.GetUser)
	rg.PUT("/:id", h.UpdateUser)
	rg.PUT("/:id/status", h.UpdateUserStatus)
	rg.POST("/:id/reset-password", h.ResetPassword)
	rg.DELETE("/:id", h.DeleteUser)
}

func (h *Handler) ListUsers(c *gin.Context) {
	var q ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	result, err := h.svc.ListUsersPage(c.Request.Context(), q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	created, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, created)
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	u, err := h.svc.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, u)
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	updated, err := h.svc.UpdateUser(c.Request.Context(), uint(id), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, updated)
}

func (h *Handler) BatchDeleteUsers(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	if err := h.svc.DeleteUserBatch(c.Request.Context(), req.IDs); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	if err := h.svc.DeleteUser(c.Request.Context(), uint(id)); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	updated, err := h.svc.UpdateUserStatus(c.Request.Context(), uint(id), req.Status)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, updated)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	if err := h.svc.ResetPassword(c.Request.Context(), uint(id), req.Password); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"reset": true})
}

func (h *Handler) ExportUsers(c *gin.Context) {
	var q ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	data, err := h.svc.ExportUsers(c.Request.Context(), q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=\"users.xlsx\"")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
