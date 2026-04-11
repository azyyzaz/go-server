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

// ListTypes godoc
//
//	@Summary		查询字典类型列表
//	@Description	按名称或编码筛选字典类型。
//	@Tags			字典管理
//	@Produce		json
//	@Param			name	query		string	false	"字典类型名称"
//	@Param			code	query		string	false	"字典类型编码"
//	@Success		200		{object}	response.Body{data=[]DictTypeResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/types [get]
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

// CreateType godoc
//
//	@Summary		创建字典类型
//	@Description	创建新的字典类型，用于归类字典数据项。
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateDictTypeRequest	true	"字典类型创建参数"
//	@Success		201		{object}	response.Body{data=DictTypeResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/types [post]
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

// GetType godoc
//
//	@Summary		获取字典类型详情
//	@Description	根据字典类型 ID 查询详细信息。
//	@Tags			字典管理
//	@Produce		json
//	@Param			id	path		int	true	"字典类型 ID"
//	@Success		200	{object}	response.Body{data=DictTypeResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/types/{id} [get]
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

// UpdateType godoc
//
//	@Summary		更新字典类型
//	@Description	修改字典类型的名称、编码、备注和状态。
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"字典类型 ID"
//	@Param			body	body		UpdateDictTypeRequest	true	"字典类型更新参数"
//	@Success		200		{object}	response.Body{data=DictTypeResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/types/{id} [put]
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

// DeleteType godoc
//
//	@Summary		删除字典类型
//	@Description	根据字典类型 ID 删除字典类型。
//	@Tags			字典管理
//	@Produce		json
//	@Param			id	path		int	true	"字典类型 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/types/{id} [delete]
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

// ListData godoc
//
//	@Summary		查询字典数据列表
//	@Description	按字典类型、标签、状态等条件筛选字典数据项。
//	@Tags			字典管理
//	@Produce		json
//	@Param			type_id		query		int		false	"字典类型 ID"
//	@Param			type_code	query		string	false	"字典类型编码"
//	@Param			label		query		string	false	"字典标签"
//	@Param			status		query		int		false	"状态：0 禁用，1 启用"
//	@Success		200			{object}	response.Body{data=[]DictDataResponse}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/items [get]
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

// CreateData godoc
//
//	@Summary		创建字典数据
//	@Description	在指定字典类型下新增一个字典数据项。
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			body	body		CreateDictDataRequest	true	"字典数据创建参数"
//	@Success		201		{object}	response.Body{data=DictDataResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/items [post]
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

// GetData godoc
//
//	@Summary		获取字典数据详情
//	@Description	根据字典数据 ID 查询详细信息。
//	@Tags			字典管理
//	@Produce		json
//	@Param			id	path		int	true	"字典数据 ID"
//	@Success		200	{object}	response.Body{data=DictDataResponse}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/items/{id} [get]
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

// UpdateData godoc
//
//	@Summary		更新字典数据
//	@Description	修改字典数据项的标签、值、排序、状态和备注。
//	@Tags			字典管理
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"字典数据 ID"
//	@Param			body	body		UpdateDictDataRequest	true	"字典数据更新参数"
//	@Success		200		{object}	response.Body{data=DictDataResponse}
//	@Failure		400		{object}	response.Body
//	@Failure		401		{object}	response.Body
//	@Failure		404		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/items/{id} [put]
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

// DeleteData godoc
//
//	@Summary		删除字典数据
//	@Description	根据字典数据 ID 删除字典数据项。
//	@Tags			字典管理
//	@Produce		json
//	@Param			id	path		int	true	"字典数据 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/items/{id} [delete]
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

// LookupByTypeCode godoc
//
//	@Summary		按类型编码查询字典数据
//	@Description	根据字典类型编码返回对应的可用字典数据项。
//	@Tags			字典管理
//	@Produce		json
//	@Param			code	path		string	true	"字典类型编码"
//	@Success		200		{object}	response.Body{data=[]DictDataResponse}
//	@Failure		401		{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/dicts/lookup/{code} [get]
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
