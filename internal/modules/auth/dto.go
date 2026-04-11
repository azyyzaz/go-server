package auth

type LoginRequest struct {
	Username    string `json:"username" binding:"required" example:"admin"`
	Password    string `json:"password" binding:"required" example:"admin123"`
	CaptchaID   string `json:"captcha_id" binding:"required" example:"captcha-8f3a2c1d"`
	CaptchaCode string `json:"captcha_code" binding:"required" example:"1234"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.access-token"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.refresh-token"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.access-token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.refresh-token"`
}
