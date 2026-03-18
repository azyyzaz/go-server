package user

import "time"

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required,min=2,max=50"`
	Email string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func toResponse(u User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
