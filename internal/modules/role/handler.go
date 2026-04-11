package role

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
	rg.GET("", h.ListRoles)
	rg.POST("", h.CreateRole)
	rg.GET("/:id", h.GetRole)
	rg.PUT("/:id", h.UpdateRole)
	rg.DELETE("/:id", h.DeleteRole)
	rg.GET("/:id/menus", h.GetRoleMenus)
	rg.PUT("/:id/menus", h.UpdateRoleMenus)
	rg.GET("/:id/apis", h.GetRoleAPIs)
	rg.PUT("/:id/apis", h.UpdateRoleAPIs)
	rg.GET("/:id/users", h.ListRoleUsers)
}

// ListRoles godoc
//
//	@Summary		分页查询角色列表
//	@Description	按角色名称、编码、状态等条件筛选角色，并返回分页结果。
//	@Tags			角色管理
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Param			name		query		string	false	"角色名称"
//	@Param			code		query		string	false	"角色编码"
//	@Param			status		query		int		false	"状态：0 禁用，1 启用"
//	@Success		200			{object}	response.Body{data=RolePageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles [get]
func (h *Handler) ListRoles(c *gin.Context) {
	var q ListRolesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	result, err := h.svc.ListRolesPage(c.Request.Context(), q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// CreateRole godoc
//
//	@Summary		创建角色
//	@Description	创建新的角色，并设置基础信息和启用状态。
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateRoleRequest	true	"角色创建参数"
//	@Success		201		{object}	response.Body{data=RoleResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	created, err := h.svc.CreateRole(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, created)
}

// GetRole godoc
//
//	@Summary		获取角色详情
//	@Description	根据角色 ID 查询角色详细信息。
//	@Tags			角色管理
//	@Produce		json
//	@Param			id	path		int	true	"角色 ID"
//	@Success		200	{object}	response.Body{data=RoleResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id} [get]
func (h *Handler) GetRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	role, err := h.svc.GetRole(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, role)
}

// UpdateRole godoc
//
//	@Summary		更新角色信息
//	@Description	修改角色名称、编码、备注和状态等信息。
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"角色 ID"
//	@Param			body	body		UpdateRoleRequest	true	"角色更新参数"
//	@Success		200		{object}	response.Body{data=RoleResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id} [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	role, err := h.svc.UpdateRole(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, role)
}

// DeleteRole godoc
//
//	@Summary		删除角色
//	@Description	根据角色 ID 删除角色。
//	@Tags			角色管理
//	@Produce		json
//	@Param			id	path		int	true	"角色 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id} [delete]
func (h *Handler) DeleteRole(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteRole(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// GetRoleMenus godoc
//
//	@Summary		获取角色菜单权限
//	@Description	查询角色已勾选的菜单 ID 及完整菜单树结构。
//	@Tags			角色管理
//	@Produce		json
//	@Param			id	path		int	true	"角色 ID"
//	@Success		200	{object}	response.Body{data=RoleMenusResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id}/menus [get]
func (h *Handler) GetRoleMenus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	result, err := h.svc.GetRoleMenus(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// UpdateRoleMenus godoc
//
//	@Summary		更新角色菜单权限
//	@Description	为角色分配菜单权限。
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"角色 ID"
//	@Param			body	body		AssignMenusRequest	true	"菜单 ID 列表"
//	@Success		200		{object}	response.Body{data=map[string]bool}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id}/menus [put]
func (h *Handler) UpdateRoleMenus(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	if err := h.svc.UpdateRoleMenus(c.Request.Context(), id, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

// GetRoleAPIs godoc
//
//	@Summary		获取角色接口权限
//	@Description	查询角色当前拥有的接口访问权限。
//	@Tags			角色管理
//	@Produce		json
//	@Param			id	path		int	true	"角色 ID"
//	@Success		200	{object}	response.Body{data=RoleAPIsResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id}/apis [get]
func (h *Handler) GetRoleAPIs(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	result, err := h.svc.GetRoleAPIs(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// UpdateRoleAPIs godoc
//
//	@Summary		更新角色接口权限
//	@Description	为角色分配接口访问权限。
//	@Tags			角色管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"角色 ID"
//	@Param			body	body		AssignAPIsRequest	true	"接口权限列表"
//	@Success		200		{object}	response.Body{data=map[string]bool}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id}/apis [put]
func (h *Handler) UpdateRoleAPIs(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req AssignAPIsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	if err := h.svc.UpdateRoleAPIs(c.Request.Context(), id, req); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

// ListRoleUsers godoc
//
//	@Summary		查询角色下的用户列表
//	@Description	分页查询绑定到指定角色的用户。
//	@Tags			角色管理
//	@Produce		json
//	@Param			id			path		int	true	"角色 ID"
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页条数"
//	@Success		200			{object}	response.Body{data=RoleUserPageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/roles/{id}/users [get]
func (h *Handler) ListRoleUsers(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var q ListRoleUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	result, err := h.svc.ListRoleUsers(c.Request.Context(), id, q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return 0, false
	}
	return uint(id), true
}
