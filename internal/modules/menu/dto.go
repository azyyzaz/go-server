package menu

import "time"

type CreateMenuRequest struct {
	ParentID   *uint  `json:"parent_id" swaggertype:"integer" example:"0"`
	Name       string `json:"name" binding:"required,min=1,max=100" example:"用户管理"`
	Type       string `json:"type" binding:"required,oneof=directory menu button" example:"menu"`
	Path       string `json:"path" example:"/system/user"`
	Component  string `json:"component" example:"/system/user/index"`
	Permission string `json:"permission" example:"system:user:list"`
	Sort       int    `json:"sort" example:"1"`
	Visible    int8   `json:"visible" binding:"oneof=0 1" example:"1"`
	Status     int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type UpdateMenuRequest struct {
	ParentID   *uint  `json:"parent_id" swaggertype:"integer" example:"0"`
	Name       string `json:"name" binding:"required,min=1,max=100" example:"用户管理"`
	Type       string `json:"type" binding:"required,oneof=directory menu button" example:"menu"`
	Path       string `json:"path" example:"/system/user"`
	Component  string `json:"component" example:"/system/user/index"`
	Permission string `json:"permission" example:"system:user:list"`
	Sort       int    `json:"sort" example:"1"`
	Visible    int8   `json:"visible" binding:"oneof=0 1" example:"1"`
	Status     int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type MenuResponse struct {
	ID         uint      `json:"id" example:"1"`
	ParentID   *uint     `json:"parent_id,omitempty" swaggertype:"integer" example:"0"`
	Name       string    `json:"name" example:"用户管理"`
	Type       string    `json:"type" example:"menu"`
	Path       string    `json:"path" example:"/system/user"`
	Component  string    `json:"component" example:"/system/user/index"`
	Permission string    `json:"permission" example:"system:user:list"`
	Sort       int       `json:"sort" example:"1"`
	Visible    int8      `json:"visible" example:"1"`
	Status     int8      `json:"status" example:"1"`
	CreatedAt  time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
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

type UpdateMenuSortsRequest struct {
	Items []MenuSortItem `json:"items" binding:"required,min=1"`
}

type MenuSortItem struct {
	ID   uint `json:"id" binding:"required" example:"1"`
	Sort int  `json:"sort" example:"10"`
}

func toResponse(item Menu) MenuResponse {
	return MenuResponse{
		ID:         item.ID,
		ParentID:   item.ParentID,
		Name:       item.Name,
		Type:       item.Type,
		Path:       item.Path,
		Component:  item.Component,
		Permission: item.Permission,
		Sort:       item.Sort,
		Visible:    item.Visible,
		Status:     item.Status,
		CreatedAt:  item.CreatedAt,
	}
}
