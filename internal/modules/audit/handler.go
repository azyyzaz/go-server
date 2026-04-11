package audit

import (
	"errors"

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
	rg.GET("/operation-logs", h.ListOperationLogs)
	rg.GET("/login-logs", h.ListLoginLogs)
}

// ListOperationLogs godoc
//
//	@Summary		查询操作日志
//	@Description	按用户名、请求方法、路径、状态码和时间范围筛选操作日志。
//	@Tags			审计日志
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Param			username	query		string	false	"用户名"
//	@Param			method		query		string	false	"请求方法"
//	@Param			path		query		string	false	"请求路径"
//	@Param			status		query		int		false	"响应状态码"
//	@Param			start_time	query		string	false	"开始时间"
//	@Param			end_time	query		string	false	"结束时间"
//	@Success		200			{object}	response.Body{data=OperationLogPageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/audits/operation-logs [get]
func (h *Handler) ListOperationLogs(c *gin.Context) {
	var q ListOperationLogsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	result, err := h.svc.ListOperationLogs(c.Request.Context(), q)
	if err != nil {
		if errors.Is(err, ErrAuditInvalidTimeRange) {
			_ = c.Error(errcode.ErrInvalidParam.AsError())
			return
		}
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}

// ListLoginLogs godoc
//
//	@Summary		查询登录日志
//	@Description	按用户名、IP、登录结果和时间范围筛选登录日志。
//	@Tags			审计日志
//	@Produce		json
//	@Param			page		query		int		false	"页码"
//	@Param			page_size	query		int		false	"每页条数"
//	@Param			username	query		string	false	"用户名"
//	@Param			ip			query		string	false	"IP 地址"
//	@Param			success		query		bool	false	"登录是否成功"
//	@Param			start_time	query		string	false	"开始时间"
//	@Param			end_time	query		string	false	"结束时间"
//	@Success		200			{object}	response.Body{data=LoginLogPageResult}
//	@Failure		400			{object}	response.Body
//	@Failure		401			{object}	response.Body
//	@Security		BearerAuth
//	@Router			/system/audits/login-logs [get]
func (h *Handler) ListLoginLogs(c *gin.Context) {
	var q ListLoginLogsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errcode.ErrInvalidParam.AsError())
		return
	}

	result, err := h.svc.ListLoginLogs(c.Request.Context(), q)
	if err != nil {
		if errors.Is(err, ErrAuditInvalidTimeRange) {
			_ = c.Error(errcode.ErrInvalidParam.AsError())
			return
		}
		_ = c.Error(err)
		return
	}
	response.Success(c, result)
}
