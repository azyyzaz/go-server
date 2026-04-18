package file

import (
	"errors"
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
	rg.POST("/upload", h.Upload)
	rg.GET("", h.List)
	rg.DELETE("/:id", h.Delete)
}

// Upload godoc
//
//	@Summary		上传文件
//	@Description	支持图片和常见办公文档上传，按日期目录存储并返回访问地址。
//	@Tags			文件管理
//	@Accept			mpfd
//	@Produce		json
//	@Param			file		formData	file	true	"上传文件"
//	@Param			category	formData	string	false	"文件分类"
//	@Success		201			{object}	response.Body{data=FileResponse}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/files/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(validation.BindError(err))
		return
	}

	var uploaderID *uint
	if raw, ok := c.Get(middleware.ContextKeyUserID); ok {
		if id, ok := raw.(uint); ok {
			uploaderID = &id
		}
	}

	result, err := h.svc.Upload(c.Request.Context(), uploaderID, fileHeader, UploadRequest{
		Category: c.PostForm("category"),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Created(c, result)
}

// List godoc
//
//	@Summary		查询文件列表
//	@Description	按关键字、分类、存储策略、上传人和时间范围分页查询文件。
//	@Tags			文件管理
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Param			keyword		query		string	false	"文件名关键字"
//	@Param			category	query		string	false	"文件分类"
//	@Param			storage		query		string	false	"存储策略"
//	@Param			uploader_id	query		int		false	"上传人 ID"
//	@Param			start_time	query		string	false	"开始时间"
//	@Param			end_time	query		string	false	"结束时间"
//	@Success		200			{object}	response.Body{data=FilePageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/files [get]
func (h *Handler) List(c *gin.Context) {
	var q ListFilesQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	result, err := h.svc.List(c.Request.Context(), q)
	if err != nil {
		if errors.Is(err, ErrFileInvalidTimeRange) {
			_ = c.Error(errcode.ErrInvalidParam.AsError())
			return
		}
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// Delete godoc
//
//	@Summary		删除文件
//	@Description	删除文件记录并清理对应存储对象。
//	@Tags			文件管理
//	@Produce		json
//	@Param			id	path		int	true	"文件 ID"
//	@Success		200	{object}	response.Body{data=map[string]bool}
//	@Failure		400	{object}	response.Body
//	@Failure		401	{object}	response.Body
//	@Failure		404	{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/files/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
