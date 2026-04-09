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

// ListUsers godoc
//
//	@Summary		用户列表（分页+搜索+筛选）
//	@Tags			Users
//	@Produce		json
//	@Param			page		query		int		false	"页码（默认1）"
//	@Param			page_size	query		int		false	"每页条数（默认10，最大100）"
//	@Param			username	query		string	false	"用户名（模糊）"
//	@Param			name		query		string	false	"姓名（模糊）"
//	@Param			status		query		int		false	"状态 1=启用 0=禁用"
//	@Success		200			{object}	response.Body{data=UserPageResult}
//	@Security		BearerAuth
//	@Router			/system/users [get]
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

// CreateUser godoc
//
//	@Summary		新增用户（含角色分配）
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateUserRequest	true	"用户信息，role_ids 为角色ID列表（可选）"
//	@Success		201		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users [post]
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

// GetUser godoc
//
//	@Summary		获取用户详情
//	@Tags			Users
//	@Produce		json
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	response.Body{data=UserResponse}
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id} [get]
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

// UpdateUser godoc
//
//	@Summary		编辑用户（含角色更新）
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"用户ID"
//	@Param			body	body		UpdateUserRequest	true	"用户信息"
//	@Success		200		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id} [put]
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

// BatchDeleteUsers godoc
//
//	@Summary		批量删除用户
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		BatchDeleteRequest	true	"用户ID列表"
//	@Success		200		{object}	response.Body
//	@Failure		400		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/batch-delete [post]
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

// DeleteUser godoc
//
//	@Summary		删除用户
//	@Tags			Users
//	@Produce		json
//	@Param			id	path		int	true	"用户ID"
//	@Success		200	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id} [delete]
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

// UpdateUserStatus godoc
//
//	@Summary		启用/禁用用户
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"用户ID"
//	@Param			body	body		UpdateStatusRequest	true	"状态 1=启用 0=禁用"
//	@Success		200		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id}/status [put]
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

// ResetPassword godoc
//
//	@Summary		重置用户密码
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"用户ID"
//	@Param			body	body		ResetPasswordRequest	true	"新密码（最少6位）"
//	@Success		200		{object}	response.Body
//	@Failure		400		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id}/reset-password [post]
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

// ExportUsers godoc
//
//	@Summary		导出用户列表（CSV）
//	@Tags			Users
//	@Produce		text/csv
//	@Param			username	query		string	false	"用户名（模糊）"
//	@Param			name		query		string	false	"姓名（模糊）"
//	@Param			status		query		int		false	"状态 1=启用 0=禁用"
//	@Success		200			{file}		binary
//	@Security		BearerAuth
//	@Router			/system/users/export [get]
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
	c.Header("Content-Disposition", "attachment; filename=\"users.csv\"")
	c.Data(200, "text/csv; charset=utf-8", data)
}
