package dept

import "time"

type CreateDeptRequest struct {
	ParentID *uint  `json:"parent_id" swaggertype:"integer" example:"0"`
	Name     string `json:"name" binding:"required,min=1,max=100" example:"研发部"`
	Leader   string `json:"leader" example:"李经理"`
	Phone    string `json:"phone" example:"13800138001"`
	Email    string `json:"email" binding:"omitempty,email" example:"rd@example.com"`
	Sort     int    `json:"sort" example:"1"`
	Status   int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type UpdateDeptRequest struct {
	ParentID *uint  `json:"parent_id" swaggertype:"integer" example:"0"`
	Name     string `json:"name" binding:"required,min=1,max=100" example:"研发部"`
	Leader   string `json:"leader" example:"李经理"`
	Phone    string `json:"phone" example:"13800138001"`
	Email    string `json:"email" binding:"omitempty,email" example:"rd@example.com"`
	Sort     int    `json:"sort" example:"1"`
	Status   int8   `json:"status" binding:"oneof=0 1" example:"1"`
}

type DeptResponse struct {
	ID        uint      `json:"id" example:"1"`
	ParentID  *uint     `json:"parent_id,omitempty" swaggertype:"integer" example:"0"`
	Name      string    `json:"name" example:"研发部"`
	Leader    string    `json:"leader" example:"李经理"`
	Phone     string    `json:"phone" example:"13800138001"`
	Email     string    `json:"email" example:"rd@example.com"`
	Sort      int       `json:"sort" example:"1"`
	Status    int8      `json:"status" example:"1"`
	CreatedAt time.Time `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

type DeptTreeNode struct {
	ID       uint           `json:"id" example:"1"`
	ParentID *uint          `json:"parent_id,omitempty" swaggertype:"integer" example:"0"`
	Name     string         `json:"name" example:"研发部"`
	Leader   string         `json:"leader" example:"李经理"`
	Phone    string         `json:"phone" example:"13800138001"`
	Email    string         `json:"email" example:"rd@example.com"`
	Sort     int            `json:"sort" example:"1"`
	Status   int8           `json:"status" example:"1"`
	Children []DeptTreeNode `json:"children,omitempty"`
}

type ListDeptUsersQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
}

type DeptUserResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"zhangsan"`
	Name     string `json:"name" example:"张三"`
	Email    string `json:"email" example:"zhangsan@example.com"`
	Phone    string `json:"phone" example:"13800138000"`
	Status   int8   `json:"status" example:"1"`
}

type DeptUserPageResult struct {
	Items    []DeptUserResponse `json:"items"`
	Total    int64              `json:"total" example:"1"`
	Page     int                `json:"page" example:"1"`
	PageSize int                `json:"page_size" example:"10"`
}

func toResponse(item Dept) DeptResponse {
	return DeptResponse{
		ID:        item.ID,
		ParentID:  item.ParentID,
		Name:      item.Name,
		Leader:    item.Leader,
		Phone:     item.Phone,
		Email:     item.Email,
		Sort:      item.Sort,
		Status:    item.Status,
		CreatedAt: item.CreatedAt,
	}
}
