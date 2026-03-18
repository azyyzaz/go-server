package user

import (
	"net/http"

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
	rg.POST("", h.CreateUser)
	rg.GET("", h.ListUsers)
	rg.GET("/:id", h.GetUser)
	rg.DELETE("/:id", h.DeleteUser)
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidArgument, "invalid request body")
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
	id := c.Param("id")
	user, err := h.svc.GetUser(c.Request.Context(), id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, user)
}

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"items": users, "total": len(users)})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteUser(c.Request.Context(), id); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
