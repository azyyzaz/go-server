package dept

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
	rg.GET("", h.ListDepts)
	rg.POST("", h.CreateDept)
	rg.GET("/:id", h.GetDept)
	rg.PUT("/:id", h.UpdateDept)
	rg.DELETE("/:id", h.DeleteDept)
	rg.GET("/:id/users", h.ListDeptUsers)
}

func (h *Handler) ListDepts(c *gin.Context) {
	items, err := h.svc.ListDepts(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, items)
}

func (h *Handler) CreateDept(c *gin.Context) {
	var req CreateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	created, err := h.svc.CreateDept(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, created)
}

func (h *Handler) GetDept(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.svc.GetDept(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) UpdateDept(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req UpdateDeptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	item, err := h.svc.UpdateDept(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) DeleteDept(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := h.svc.DeleteDept(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *Handler) ListDeptUsers(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var q ListDeptUsersQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	result, err := h.svc.ListDeptUsers(c.Request.Context(), id, q)
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
