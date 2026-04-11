package dict

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
	rg.GET("/types", h.ListTypes)
	rg.POST("/types", h.CreateType)
	rg.GET("/types/:id", h.GetType)
	rg.PUT("/types/:id", h.UpdateType)
	rg.DELETE("/types/:id", h.DeleteType)

	rg.GET("/items", h.ListData)
	rg.POST("/items", h.CreateData)
	rg.GET("/items/:id", h.GetData)
	rg.PUT("/items/:id", h.UpdateData)
	rg.DELETE("/items/:id", h.DeleteData)

	rg.GET("/lookup/:code", h.LookupByTypeCode)
}

func (h *Handler) ListTypes(c *gin.Context) {
	var q ListDictTypesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	items, err := h.svc.ListTypes(c.Request.Context(), q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, items)
}

func (h *Handler) CreateType(c *gin.Context) {
	var req CreateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	item, err := h.svc.CreateType(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) GetType(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := h.svc.GetType(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) UpdateType(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req UpdateDictTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	item, err := h.svc.UpdateType(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) DeleteType(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteType(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *Handler) ListData(c *gin.Context) {
	var q ListDictDataQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	items, err := h.svc.ListData(c.Request.Context(), q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, items)
}

func (h *Handler) CreateData(c *gin.Context) {
	var req CreateDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	item, err := h.svc.CreateData(c.Request.Context(), req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, item)
}

func (h *Handler) GetData(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	item, err := h.svc.GetData(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) UpdateData(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	var req UpdateDictDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	item, err := h.svc.UpdateData(c.Request.Context(), id, req)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, item)
}

func (h *Handler) DeleteData(c *gin.Context) {
	id, ok := parseID(c, "id")
	if !ok {
		return
	}
	if err := h.svc.DeleteData(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *Handler) LookupByTypeCode(c *gin.Context) {
	items, err := h.svc.LookupByTypeCode(c.Request.Context(), c.Param("code"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, items)
}

func parseID(c *gin.Context, key string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return 0, false
	}
	return uint(id), true
}
