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
