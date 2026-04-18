package user

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50" example:"zhangsan"`
	Password string `json:"password" binding:"required,min=6" example:"Admin123!"`
	Name     string `json:"name" binding:"required,min=2,max=50" example:"Zhang San"`
	Email    string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Phone    string `json:"phone" example:"13800138000"`
	DeptID   *uint  `json:"dept_id" swaggertype:"integer" example:"1"`
	RoleIDs  []uint `json:"role_ids" swaggertype:"array,integer" example:"1,2"`
}

type UpdateUserRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=50" example:"Zhang San"`
	Email   string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Phone   string `json:"phone" example:"13800138000"`
	DeptID  *uint  `json:"dept_id" swaggertype:"integer" example:"1"`
	RoleIDs []uint `json:"role_ids" swaggertype:"array,integer" example:"1,2"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50" example:"Zhang San"`
	Email string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Phone string `json:"phone" example:"13800138000"`
}

type ListUsersQuery struct {
	Page     int    `form:"page" binding:"omitempty,min=1" default:"1" example:"1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100" default:"10" example:"10"`
	Username string `form:"username" example:"zhangsan"`
	Name     string `form:"name" example:"Zhang San"`
	Status   *int8  `form:"status" binding:"omitempty,oneof=0 1" swaggertype:"integer" example:"1"`
}

type RoleInfo struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Administrator"`
	Code string `json:"code" example:"admin"`
}

type UserResponse struct {
	ID        uint       `json:"id" example:"1"`
	Username  string     `json:"username" example:"zhangsan"`
	Name      string     `json:"name" example:"Zhang San"`
	Email     string     `json:"email" example:"zhangsan@example.com"`
	Phone     string     `json:"phone" example:"13800138000"`
	Avatar    string     `json:"avatar" example:"/uploads/avatars/1_1713072000.png"`
	DeptID    *uint      `json:"dept_id,omitempty" swaggertype:"integer" example:"1"`
	Status    int8       `json:"status" example:"1"`
	Roles     []RoleInfo `json:"roles"`
	CreatedAt time.Time  `json:"created_at" swaggertype:"string" example:"2026-04-11T12:00:00Z"`
}

type UserPageResult struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total" example:"1"`
	Page     int            `json:"page" example:"1"`
	PageSize int            `json:"page_size" example:"10"`
}

type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1" swaggertype:"array,integer" example:"1,2,3"`
}

type UpdateStatusRequest struct {
	Status int8 `json:"status" binding:"oneof=0 1" example:"1"`
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6" example:"NewPass123!"`
}

func toResponse(u User) UserResponse {
	roles := make([]RoleInfo, 0, len(u.Roles))
	for _, r := range u.Roles {
		roles = append(roles, RoleInfo{ID: r.ID, Name: r.Name, Code: r.Code})
	}
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		DeptID:    u.DeptID,
		Status:    u.Status,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
	}
}
