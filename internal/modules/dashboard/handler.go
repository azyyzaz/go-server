package dashboard

import (
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
	rg.GET("/overview", h.GetOverview)
	rg.GET("/sales-trend", h.GetSalesTrend)
	rg.GET("/visit-sources", h.GetVisitSources)
	rg.GET("/sales-categories", h.GetSalesCategories)
	rg.GET("/online-users", h.GetOnlineUsers)
}

func (h *Handler) GetOverview(c *gin.Context) {
	result, err := h.svc.GetOverview(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) GetSalesTrend(c *gin.Context) {
	var q SalesTrendQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	result, err := h.svc.GetSalesTrend(c.Request.Context(), q)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) GetVisitSources(c *gin.Context) {
	result, err := h.svc.GetVisitSources(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) GetSalesCategories(c *gin.Context) {
	result, err := h.svc.GetSalesCategories(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

func (h *Handler) GetOnlineUsers(c *gin.Context) {
	result, err := h.svc.GetOnlineUsers(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}
