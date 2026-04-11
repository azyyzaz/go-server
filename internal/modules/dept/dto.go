package dept

import "time"

type CreateDeptRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status" binding:"oneof=0 1"`
}

type UpdateDeptRequest struct {
	ParentID *uint  `json:"parent_id"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Leader   string `json:"leader"`
	Phone    string `json:"phone"`
	Email    string `json:"email" binding:"omitempty,email"`
	Sort     int    `json:"sort"`
	Status   int8   `json:"status" binding:"oneof=0 1"`
}

type DeptResponse struct {
	ID        uint      `json:"id"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	Name      string    `json:"name"`
	Leader    string    `json:"leader"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Sort      int       `json:"sort"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DeptTreeNode struct {
	ID       uint           `json:"id"`
	ParentID *uint          `json:"parent_id,omitempty"`
	Name     string         `json:"name"`
	Leader   string         `json:"leader"`
	Phone    string         `json:"phone"`
	Email    string         `json:"email"`
	Sort     int            `json:"sort"`
	Status   int8           `json:"status"`
	Children []DeptTreeNode `json:"children,omitempty"`
}

type ListDeptUsersQuery struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100"`
}

type DeptUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
}

type DeptUserPageResult struct {
	Items    []DeptUserResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
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
