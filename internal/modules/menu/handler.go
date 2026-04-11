package menu

import (
	"strconv"

	"go-server/internal/errcode"
	"go-server/internal/middleware"
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
	rg.GET("", h.ListMenus)
	rg.GET("/current", h.ListCurrentUserMenus)
	rg.PUT("/sort", h.UpdateMenuSorts)
	rg.POST("", h.CreateMenu)
	rg.GET("/:id", h.GetMenu)
	rg.PUT("/:id", h.UpdateMenu)
	rg.DELETE("/:id", h.DeleteMenu)
}

func (h *Handler) ListMenus(c *gin.Context) {
	result, err := h.svc.ListMenus(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

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

func (h *Handler) CreateMenu(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	created, err := h.svc.CreateMenu(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, created)
}

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

func (h *Handler) UpdateMenu(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	item, err := h.svc.UpdateMenu(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

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

func (h *Handler) UpdateMenuSorts(c *gin.Context) {
	var req UpdateMenuSortsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
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
