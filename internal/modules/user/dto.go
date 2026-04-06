package user

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"     binding:"required,min=2,max=50"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"`
	RoleIDs  []uint `json:"role_ids"`
}

type UpdateUserRequest struct {
	Name    string `json:"name"  binding:"required,min=2,max=50"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone"`
	RoleIDs []uint `json:"role_ids"`
}

type ListUsersQuery struct {
	Page     int    `form:"page"      binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Username string `form:"username"`
	Name     string `form:"name"`
	Status   *int8  `form:"status"`
}

type RoleInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Status    int8       `json:"status"`
	Roles     []RoleInfo `json:"roles"`
	CreatedAt time.Time  `json:"created_at"`
}

type UserPageResult struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
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
		Status:    u.Status,
		Roles:     roles,
		CreatedAt: u.CreatedAt,
	}
}
