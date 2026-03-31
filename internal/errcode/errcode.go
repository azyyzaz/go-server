package errcode

import (
	"net/http"

	"go-server/internal/response"
)

// ErrCode 是业务错误码类型（整数，便于枚举和映射）
type ErrCode int

// 每个 ErrCode 携带的元数据
type meta struct {
	status int    // HTTP 状态码
	code   string // 响应 body 里的 code 字段
	msgEN  string // 英文 message（给日志）
	msgZH  string // 中文 message（给前端）
}

// AsError 将错误码转成 AppError，可直接 return
func (e ErrCode) AsError() *response.AppError {
	m := table[e]
	return response.NewAppError(m.status, m.code, m.msgZH)
}

// HTTP 状态码
func (e ErrCode) Status() int { return table[e].status }

// ── 错误码定义 ───────────────────────────────────────────

const (
	// 通用
	ErrInvalidParam  ErrCode = 10001
	ErrUnauthorized  ErrCode = 10002
	ErrForbidden     ErrCode = 10003
	ErrNotFound      ErrCode = 10004
	ErrInternalError ErrCode = 10005

	// 用户模块  2xxxx
	ErrUserNotFound    ErrCode = 20001
	ErrUserEmailExists ErrCode = 20002
)

// ── 元数据映射表 ─────────────────────────────────────────

var table = map[ErrCode]meta{
	ErrInvalidParam:  {http.StatusBadRequest, "INVALID_ARGUMENT", "invalid param", "参数错误"},
	ErrUnauthorized:  {http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized", "未登录或 Token 已过期"},
	ErrForbidden:     {http.StatusForbidden, "FORBIDDEN", "forbidden", "没有操作权限"},
	ErrNotFound:      {http.StatusNotFound, "NOT_FOUND", "not found", "资源不存在"},
	ErrInternalError: {http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", "服务器内部错误"},

	ErrUserNotFound:    {http.StatusNotFound, "USER_NOT_FOUND", "user not found", "用户不存在"},
	ErrUserEmailExists: {http.StatusConflict, "USER_EMAIL_EXISTS", "email already exists", "邮箱已被注册"},
}
