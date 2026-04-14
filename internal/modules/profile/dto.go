package profile

import (
	"go-server/internal/modules/audit"
	"go-server/internal/modules/user"
)

type UpdateProfileRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50" example:"寮犱笁"`
	Email string `json:"email" binding:"required,email" example:"zhangsan@example.com"`
	Phone string `json:"phone" example:"13800138000"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6" example:"Admin123!"`
	NewPassword string `json:"new_password" binding:"required,min=6" example:"NewPass123!"`
}

type AvatarUploadResponse struct {
	Avatar string `json:"avatar" example:"/uploads/avatars/1_1713072000.png"`
}

type LoginLogPageResult = audit.LoginLogPageResult
type ProfileResponse = user.UserResponse
