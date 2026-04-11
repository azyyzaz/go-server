package menu

import "time"

type CreateMenuRequest struct {
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name" binding:"required,min=1,max=100"`
	Type       string `json:"type" binding:"required,oneof=directory menu button"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Sort       int    `json:"sort"`
	Visible    int8   `json:"visible" binding:"oneof=0 1"`
	Status     int8   `json:"status" binding:"oneof=0 1"`
}

type UpdateMenuRequest struct {
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name" binding:"required,min=1,max=100"`
	Type       string `json:"type" binding:"required,oneof=directory menu button"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Permission string `json:"permission"`
	Sort       int    `json:"sort"`
	Visible    int8   `json:"visible" binding:"oneof=0 1"`
	Status     int8   `json:"status" binding:"oneof=0 1"`
}

type MenuResponse struct {
	ID         uint      `json:"id"`
	ParentID   *uint     `json:"parent_id,omitempty"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Path       string    `json:"path"`
	Component  string    `json:"component"`
	Permission string    `json:"permission"`
	Sort       int       `json:"sort"`
	Visible    int8      `json:"visible"`
	Status     int8      `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
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

type UpdateMenuSortsRequest struct {
	Items []MenuSortItem `json:"items" binding:"required,min=1"`
}

type MenuSortItem struct {
	ID   uint `json:"id" binding:"required"`
	Sort int  `json:"sort"`
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
