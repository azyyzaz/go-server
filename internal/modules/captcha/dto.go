package captcha

type GenerateResponse struct {
	CaptchaID   string `json:"captcha_id" example:"captcha-8f3a2c1d"`
	ImageBase64 string `json:"image_base64" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."`
}

type VerifyRequest struct {
	CaptchaID string `json:"captcha_id" binding:"required" example:"captcha-8f3a2c1d"`
	Answer    string `json:"answer" binding:"required" example:"1234"`
}
