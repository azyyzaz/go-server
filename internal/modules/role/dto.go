package role

import (
	"strings"
	"time"
)

type CreateRoleRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=50"`
	Code   string `json:"code" binding:"required,min=2,max=50"`
	Remark string `json:"remark"`
	Status int8   `json:"status" binding:"oneof=0 1"`
}

type UpdateRoleRequest struct {
	Name   string `json:"name" binding:"required,min=2,max=50"`
	Code   string `json:"code"`
	Remark string `json:"remark"`
	Status int8   `json:"status" binding:"oneof=0 1"`
}

type ListRolesQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1"`
}

type RoleResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Remark    string    `json:"remark"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePageResult struct {
	Items    []RoleResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type AssignMenusRequest struct {
	MenuIDs []uint `json:"menu_ids"`
}

type MenuTreeNode struct {
	ID         uint           `json:"id"`
	ParentID   *uint          `json:"parent_id,omitempty"`
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Path       string         `json:"path"`
	Component  string         `json:"component"`
	Permission string         `json:"permission"`
	Sort       int            `json:"sort"`
	Visible    int8           `json:"visible"`
	Status     int8           `json:"status"`
	Children   []MenuTreeNode `json:"children,omitempty"`
}

type RoleMenusResponse struct {
	CheckedIDs []uint         `json:"checked_ids"`
	Menus      []MenuTreeNode `json:"menus"`
}

type APIPermission struct {
	Path   string `json:"path" binding:"required"`
	Method string `json:"method" binding:"required"`
}

type AssignAPIsRequest struct {
	Permissions []APIPermission `json:"permissions"`
}

type RoleAPIsResponse struct {
	Permissions []APIPermission `json:"permissions"`
}

type ListRoleUsersQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type RoleUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
}

type RoleUserPageResult struct {
	Items    []RoleUserResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
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
