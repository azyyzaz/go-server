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

// ListDepts godoc
//
//	@Summary		查询部门树
//	@Description	返回完整部门树结构，常用于组织架构管理页面。
//	@Tags			部门管理
//	@Produce		json
//	@Success		200	{object}	response.Body{data=[]DeptTreeNode}
//	@Failure		401	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/depts [get]
func (h *Handler) ListDepts(c *gin.Context) {
	items, err := h.svc.ListDepts(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, items)
}

// CreateDept godoc
//
//	@Summary		创建部门
//	@Description	创建新的部门节点，可指定上级部门。
//	@Tags			部门管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateDeptRequest	true	"部门创建参数"
//	@Success		201		{object}	response.Body{data=DeptResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/depts [post]
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

// GetDept godoc
//
//	@Summary		获取部门详情
//	@Description	根据部门 ID 查询部门详细信息。
//	@Tags			部门管理
//	@Produce		json
//	@Param			id	path		int	true	"部门 ID"
//	@Success		200	{object}	response.Body{data=DeptResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/depts/{id} [get]
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

// UpdateDept godoc
//
//	@Summary		更新部门
//	@Description	修改部门名称、负责人、联系方式、排序和状态等信息。
//	@Tags			部门管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"部门 ID"
//	@Param			body	body		UpdateDeptRequest	true	"部门更新参数"
//	@Success		200		{object}	response.Body{data=DeptResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/depts/{id} [put]
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

// DeleteDept godoc
//
//	@Summary		删除部门
//	@Description	根据部门 ID 删除部门节点。
//	@Tags			部门管理
//	@Produce		json
//	@Param			id	path		int	true	"部门 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/depts/{id} [delete]
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

// ListDeptUsers godoc
//
//	@Summary		查询部门下的用户列表
//	@Description	分页查询归属于指定部门的用户。
//	@Tags			部门管理
//	@Produce		json
//	@Param			id			path		int	true	"部门 ID"
//	@Param			page		query		int	false	"页码"
//	@Param			page_size	query		int	false	"每页条数"
//	@Success		200			{object}	response.Body{data=DeptUserPageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/depts/{id}/users [get]
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
