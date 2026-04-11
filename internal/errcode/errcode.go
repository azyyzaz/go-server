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

	ErrMenuNotFound      ErrCode = 40001
	ErrMenuNameExists    ErrCode = 40002
	ErrMenuHasChildren   ErrCode = 40003
	ErrMenuParentInvalid ErrCode = 40004

	ErrDeptNotFound      ErrCode = 50001
	ErrDeptNameExists    ErrCode = 50002
	ErrDeptHasChildren   ErrCode = 50003
	ErrDeptHasUsers      ErrCode = 50004
	ErrDeptParentInvalid ErrCode = 50005

	ErrDictTypeNotFound    ErrCode = 60001
	ErrDictTypeCodeExists  ErrCode = 60002
	ErrDictTypeHasItems    ErrCode = 60003
	ErrDictDataNotFound    ErrCode = 60004
	ErrDictDataLabelExists ErrCode = 60005
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

	ErrMenuNotFound:      {http.StatusNotFound, "MENU_NOT_FOUND", "menu not found", "菜单不存在"},
	ErrMenuNameExists:    {http.StatusConflict, "MENU_NAME_EXISTS", "menu name already exists", "同级菜单名称已存在"},
	ErrMenuHasChildren:   {http.StatusConflict, "MENU_HAS_CHILDREN", "menu still has children", "该菜单下仍有子节点，不能删除"},
	ErrMenuParentInvalid: {http.StatusBadRequest, "MENU_PARENT_INVALID", "invalid menu parent", "菜单父节点不合法"},

	ErrDeptNotFound:      {http.StatusNotFound, "DEPT_NOT_FOUND", "department not found", "部门不存在"},
	ErrDeptNameExists:    {http.StatusConflict, "DEPT_NAME_EXISTS", "department name already exists", "同级部门名称已存在"},
	ErrDeptHasChildren:   {http.StatusConflict, "DEPT_HAS_CHILDREN", "department still has children", "该部门下仍有子部门，不能删除"},
	ErrDeptHasUsers:      {http.StatusConflict, "DEPT_HAS_USERS", "department still has users", "该部门下仍有用户，不能删除"},
	ErrDeptParentInvalid: {http.StatusBadRequest, "DEPT_PARENT_INVALID", "invalid department parent", "部门父节点不合法"},

	ErrDictTypeNotFound:    {http.StatusNotFound, "DICT_TYPE_NOT_FOUND", "dict type not found", "字典类型不存在"},
	ErrDictTypeCodeExists:  {http.StatusConflict, "DICT_TYPE_CODE_EXISTS", "dict type code already exists", "字典类型标识已存在"},
	ErrDictTypeHasItems:    {http.StatusConflict, "DICT_TYPE_HAS_ITEMS", "dict type still has items", "字典类型下仍有数据，不能删除"},
	ErrDictDataNotFound:    {http.StatusNotFound, "DICT_DATA_NOT_FOUND", "dict data not found", "字典数据不存在"},
	ErrDictDataLabelExists: {http.StatusConflict, "DICT_DATA_LABEL_EXISTS", "dict data label already exists", "同一字典类型下的标签已存在"},
}
