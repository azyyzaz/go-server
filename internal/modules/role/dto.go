package role

import (
	"strings"
	"time"
)

type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=50" example:"管理员"`
	Code   string `json:"code" binding:"required,min=2,max=50" example:"admin"`
	Remark string `json:"remark" example:"系统管理员角色"`
	Status int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type UpdateRoleRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=50" example:"管理员"`
	Code   string `json:"code" example:"admin"`
	Remark string `json:"remark" example:"系统管理员角色"`
	Status int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type ListRolesQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Name     string `form:"name" example:"管理员"`
	Code     string `form:"code" example:"admin"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1" swaggertype:"integer" example:"1"`
}

type RoleResponse struct {
	ID        uint      `json:"id" example:"1"`
	Name      string    `json:"name" example:"管理员"`
	Code      string    `json:"code" example:"admin"`
	Remark    string    `json:"remark" example:"系统管理员角色"`
	Status    int8      `json:"status" example:"1"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

type RolePageResult struct {
	Items    []RoleResponse `json:"items"`
	Total    int64          `json:"total" example:"1"`
	Page     int            `json:"page" example:"1"`
	PageSize int            `json:"page_size" example:"10"`
}

type AssignMenusRequest struct {
	MenuIDs []uint `json:"menu_ids" swaggertype:"array,integer" example:"1,2,3"`
}

type MenuTreeNode struct {
	ID         uint           `json:"id" example:"1"`
	ParentID   *uint          `json:"parent_id,omitempty" swaggertype:"integer" example:"0"`
	Name       string         `json:"name" example:"系统管理"`
	Type       string         `json:"type" example:"directory"`
	Path       string         `json:"path" example:"/system"`
	Component  string         `json:"component" example:"/system/index"`
	Permission string         `json:"permission" example:"system:user:list"`
	Sort       int            `json:"sort" example:"1"`
	Visible    int8           `json:"visible" example:"1"`
	Status     int8           `json:"status" example:"1"`
	Children   []MenuTreeNode `json:"children,omitempty"`
}

type RoleMenusResponse struct {
	CheckedIDs []uint         `json:"checked_ids" swaggertype:"array,integer" example:"1,2,3"`
	Menus      []MenuTreeNode `json:"menus"`
}

type APIPermission struct {
	Path   string `json:"path" binding:"required" example:"/api/v1/system/users"`
	Method string `json:"method" binding:"required" example:"GET"`
}

type AssignAPIsRequest struct {
	Permissions []APIPermission `json:"permissions"`
}

type RoleAPIsResponse struct {
	Permissions []APIPermission `json:"permissions"`
}

type ListRoleUsersQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
}

type RoleUserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"zhangsan"`
	Name     string `json:"name" example:"张三"`
	Email    string `json:"email" example:"zhangsan@example.com"`
	Phone    string `json:"phone" example:"13800138000"`
	Status   int8   `json:"status" example:"1"`
}

type RoleUserPageResult struct {
	Items    []RoleUserResponse `json:"items"`
	Total    int64              `json:"total" example:"1"`
	Page     int                `json:"page" example:"1"`
	PageSize int                `json:"page_size" example:"10"`
}

func toResponse(r Role) RoleResponse {
	return RoleResponse{
		ID:        r.ID,
		Name:      r.Name,
		Code:      r.Code,
		Remark:    r.Remark,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
	}
}

func normalizePermission(path, method string) APIPermission {
	path = strings.TrimSpace(path)
	path = strings.TrimRight(path, "/")
	if path == "" {
		path = "/"
	}
	return APIPermission{
		Path:   path,
		Method: strings.ToUpper(strings.TrimSpace(method)),
	}
}
