package captcha

import (
	"github.com/mojocn/base64Captcha"
)

// 使用内存存储，生产环境可替换为 Redis store
var store = base64Captcha.DefaultMemStore

type Service interface {
	Generate() (GenerateResponse, error)
	Verify(id, answer string) bool
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Generate() (GenerateResponse, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 6, 0.7, 80)
	c := base64Captcha.NewCaptcha(driver, store)

	id, b64, _, err := c.Generate()
	if err != nil {
		return GenerateResponse{}, err
	}
	return GenerateResponse{CaptchaID: id, ImageBase64: b64}, nil
}

func (s *service) Verify(id, answer string) bool {
	return store.Verify(id, answer, true)
}
