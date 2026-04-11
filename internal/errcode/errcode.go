package errcode

import (
	"net/http"

	"go-server/internal/response"
)

// ErrCode represents business-level error codes.
type ErrCode int

type meta struct {
	status int
	code   string
	msgEN  string
	msgZH  string
}

func (e ErrCode) AsError() *response.AppError {
	m := table[e]
	return response.NewAppError(m.status, m.code, m.msgZH)
}

func (e ErrCode) Status() int { return table[e].status }

const (
	ErrInvalidParam  ErrCode = 10001
	ErrUnauthorized  ErrCode = 10002
	ErrForbidden     ErrCode = 10003
	ErrNotFound      ErrCode = 10004
	ErrInternalError ErrCode = 10005

	ErrUserNotFound       ErrCode = 20001
	ErrUserEmailExists    ErrCode = 20002
	ErrInvalidCredentials ErrCode = 20003
	ErrCaptchaInvalid     ErrCode = 20004
	ErrUserDisabled       ErrCode = 20005

	ErrRoleNotFound      ErrCode = 30001
	ErrRoleCodeExists    ErrCode = 30002
	ErrRoleHasUsers      ErrCode = 30003
	ErrRoleCodeImmutable ErrCode = 30004
)

var table = map[ErrCode]meta{
	ErrInvalidParam:  {http.StatusBadRequest, "INVALID_ARGUMENT", "invalid param", "参数错误"},
	ErrUnauthorized:  {http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized", "未登录或 Token 已过期"},
	ErrForbidden:     {http.StatusForbidden, "FORBIDDEN", "forbidden", "没有操作权限"},
	ErrNotFound:      {http.StatusNotFound, "NOT_FOUND", "not found", "资源不存在"},
	ErrInternalError: {http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", "服务器内部错误"},

	ErrUserNotFound:       {http.StatusNotFound, "USER_NOT_FOUND", "user not found", "用户不存在"},
	ErrUserEmailExists:    {http.StatusConflict, "USER_EMAIL_EXISTS", "email already exists", "邮箱已被注册"},
	ErrInvalidCredentials: {http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid username or password", "用户名或密码错误"},
	ErrCaptchaInvalid:     {http.StatusBadRequest, "CAPTCHA_INVALID", "captcha invalid or expired", "验证码错误或已过期"},
	ErrUserDisabled:       {http.StatusForbidden, "USER_DISABLED", "user is disabled", "用户已被禁用"},

	ErrRoleNotFound:      {http.StatusNotFound, "ROLE_NOT_FOUND", "role not found", "角色不存在"},
	ErrRoleCodeExists:    {http.StatusConflict, "ROLE_CODE_EXISTS", "role code already exists", "角色标识已存在"},
	ErrRoleHasUsers:      {http.StatusConflict, "ROLE_HAS_USERS", "role still has users", "角色下仍有用户，不能删除"},
	ErrRoleCodeImmutable: {http.StatusBadRequest, "ROLE_CODE_IMMUTABLE", "role code cannot be changed", "角色标识创建后不可修改"},
}
