package auth

import (
	"context"
	"testing"
	"time"

	"go-server/internal/jwt"
	"go-server/internal/modules/captcha"
	"go-server/internal/modules/user"
	"go-server/internal/response"
)

type captchaStub struct {
	ok bool
}

func (c captchaStub) Generate() (captcha.GenerateResponse, error) {
	return captcha.GenerateResponse{}, nil
}

func (c captchaStub) Verify(_, _ string) bool {
	return c.ok
}

func TestLoginReturnsDisabledErrorForDisabledUser(t *testing.T) {
	repo := user.NewInMemoryRepository()
	userSvc := user.NewService(repo)

	created, err := userSvc.CreateUser(context.Background(), user.CreateUserRequest{
		Username: "disabled-user",
		Password: "secret123",
		Name:     "Disabled User",
		Email:    "disabled@example.com",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	if _, err := userSvc.UpdateUserStatus(context.Background(), created.ID, 0); err != nil {
		t.Fatalf("disable user: %v", err)
	}

	svc := NewService(
		userSvc,
		jwt.NewManager("test-secret", 15*time.Minute, 24*time.Hour),
		nil,
		captchaStub{ok: true},
	)

	_, err = svc.Login(context.Background(), LoginRequest{
		Username:    "disabled-user",
		Password:    "secret123",
		CaptchaID:   "captcha-id",
		CaptchaCode: "123456",
	})
	if err == nil {
		t.Fatal("expected disabled user login to fail")
	}

	appErr, ok := err.(*response.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != "USER_DISABLED" {
		t.Fatalf("expected USER_DISABLED, got %s", appErr.Code)
	}
}
