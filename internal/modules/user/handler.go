package user

import (
	"strconv"

	"go-server/internal/errcode"
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
//	@Summary		分页查询用户列表
//	@Description	按用户名、姓名、状态等条件筛选用户，并返回分页结果。
//	@Tags			用户管理
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Param			username	query		string	false	"用户名"
//	@Param			name		query		string	false	"姓名"
//	@Param			status		query		int		false	"状态：0 禁用，1 启用"
//	@Success		200			{object}	response.Body{data=UserPageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	var q ListUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(validation.BindError(err))
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
//	@Summary		创建用户
//	@Description	创建一个新的系统用户，并绑定部门和角色信息。
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateUserRequest	true	"用户创建参数"
//	@Success		201		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
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
//	@Description	根据用户 ID 查询单个用户的详细信息。
//	@Tags			用户管理
//	@Produce		json
//	@Param			id	path		int	true	"用户 ID"
//	@Success		200	{object}	response.Body{data=UserResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(validation.BindError(err))
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
//	@Summary		更新用户信息
//	@Description	修改用户基础资料、所属部门和角色信息。
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"用户 ID"
//	@Param			body	body		UpdateUserRequest	true	"用户更新参数"
//	@Success		200		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
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
//	@Description	根据用户 ID 列表一次性删除多个用户。
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		BatchDeleteRequest	true	"待删除的用户 ID 列表"
//	@Success		200		{object}	response.Body{data=map[string]bool}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/batch-delete [post]
func (h *Handler) BatchDeleteUsers(c *gin.Context) {
	var req BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
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
//	@Description	根据用户 ID 删除单个用户。
//	@Tags			用户管理
//	@Produce		json
//	@Param			id	path		int	true	"用户 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(validation.BindError(err))
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
//	@Summary		修改用户状态
//	@Description	启用或禁用指定用户账号。
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"用户 ID"
//	@Param			body	body		UpdateStatusRequest	true	"状态更新参数"
//	@Success		200		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
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
//	@Description	为指定用户重置登录密码。
//	@Tags			用户管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"用户 ID"
//	@Param			body	body		ResetPasswordRequest	true	"密码重置参数"
//	@Success		200		{object}	response.Body{data=map[string]bool}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
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
//	@Summary		导出用户列表
//	@Description	按照筛选条件导出用户数据，返回 Excel 文件。
//	@Tags			用户管理
//	@Produce		octet-stream
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Param			username	query		string	false	"用户名"
//	@Param			name		query		string	false	"姓名"
//	@Param			status		query		int		false	"状态：0 禁用，1 启用"
//	@Success		200			{file}		file
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
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
	c.Header("Content-Disposition", "attachment; filename=\"users.xlsx\"")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}
