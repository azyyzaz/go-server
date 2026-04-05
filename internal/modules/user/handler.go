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
	rg.POST("", h.CreateUser)
	rg.GET("", h.ListUsers)
	rg.GET("/:id", h.GetUser)
	rg.DELETE("/:id", h.DeleteUser)
}

// CreateUser godoc
//
//	@Summary		创建用户
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateUserRequest	true	"用户信息"
//	@Success		201		{object}	response.Body{data=UserResponse}
//	@Failure		400		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/users [post]
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
//	@Router			/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	user, err := h.svc.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, user)
}

// ListUsers godoc
//
//	@Summary		用户列表
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	response.Body{data=[]UserResponse}
//	@Security		BearerAuth
//	@Router			/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"items": users, "total": len(users)})
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
//	@Router			/users/{id} [delete]
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
