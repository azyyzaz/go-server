package captcha

type GenerateResponse struct {
	CaptchaID string `json:"captcha_id"`
	ImageBase64 string `json:"image_base64"`
}

type VerifyRequest struct {
	CaptchaID string `json:"captcha_id" binding:"required"`
	Answer    string `json:"answer" binding:"required"`
}
