package menu

import (
	"strconv"

	"go-server/internal/errcode"
	"go-server/internal/middleware"
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
	rg.GET("", h.ListMenus)
	rg.GET("/current", h.ListCurrentUserMenus)
	rg.PUT("/sort", h.UpdateMenuSorts)
	rg.POST("", h.CreateMenu)
	rg.GET("/:id", h.GetMenu)
	rg.PUT("/:id", h.UpdateMenu)
	rg.DELETE("/:id", h.DeleteMenu)
}

// ListMenus godoc
//
//	@Summary		查询菜单树
//	@Description	返回完整菜单树，常用于菜单管理页面展示。
//	@Tags			菜单管理
//	@Produce		json
//	@Success		200	{object}	response.Body{data=[]MenuTreeNode}
//	@Failure		401	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus [get]
func (h *Handler) ListMenus(c *gin.Context) {
	result, err := h.svc.ListMenus(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListCurrentUserMenus godoc
//
//	@Summary		查询当前用户菜单
//	@Description	根据当前登录用户的权限返回可见菜单树。
//	@Tags			菜单管理
//	@Produce		json
//	@Success		200	{object}	response.Body{data=[]MenuTreeNode}
//	@Failure		401	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus/current [get]
func (h *Handler) ListCurrentUserMenus(c *gin.Context) {
	userIDValue, ok := c.Get(middleware.ContextKeyUserID)
	if !ok {
		_ = c.Error(errcode.ErrUnauthorized.AsError())
		return
	}
	userID, ok := userIDValue.(uint)
	if !ok {
		_ = c.Error(errcode.ErrUnauthorized.AsError())
		return
	}

	result, err := h.svc.ListCurrentUserMenus(c.Request.Context(), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// CreateMenu godoc
//
//	@Summary		创建菜单
//	@Description	创建目录、菜单或按钮类型的权限节点。
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateMenuRequest	true	"菜单创建参数"
//	@Success		201		{object}	response.Body{data=MenuResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus [post]
func (h *Handler) CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}
	created, err := h.svc.CreateMenu(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, created)
}

// GetMenu godoc
//
//	@Summary		获取菜单详情
//	@Description	根据菜单 ID 查询菜单详情。
//	@Tags			菜单管理
//	@Produce		json
//	@Param			id	path		int	true	"菜单 ID"
//	@Success		200	{object}	response.Body{data=MenuResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus/{id} [get]
func (h *Handler) GetMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.svc.GetMenu(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

// UpdateMenu godoc
//
//	@Summary		更新菜单
//	@Description	修改菜单节点的名称、路径、组件、权限标识和显示状态等信息。
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"菜单 ID"
//	@Param			body	body		UpdateMenuRequest	true	"菜单更新参数"
//	@Success		200		{object}	response.Body{data=MenuResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus/{id} [put]
func (h *Handler) UpdateMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}
	item, err := h.svc.UpdateMenu(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

// DeleteMenu godoc
//
//	@Summary		删除菜单
//	@Description	根据菜单 ID 删除菜单节点。
//	@Tags			菜单管理
//	@Produce		json
//	@Param			id	path		int	true	"菜单 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus/{id} [delete]
func (h *Handler) DeleteMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteMenu(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// UpdateMenuSorts godoc
//
//	@Summary		更新菜单排序
//	@Description	批量更新多个菜单节点的排序值。
//	@Tags			菜单管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		UpdateMenuSortsRequest	true	"菜单排序更新参数"
//	@Success		200		{object}	response.Body{data=map[string]bool}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/menus/sort [put]
func (h *Handler) UpdateMenuSorts(c *gin.Context) {
	var req UpdateMenuSortsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}
	if err := h.svc.UpdateMenuSorts(c.Request.Context(), req); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return 0, false
	}
	return uint(id), true
}
