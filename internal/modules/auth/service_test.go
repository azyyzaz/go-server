package auth

import (
	"context"
	"testing"
	"time"

	"go-server/internal/jwt"
	"go-server/internal/modules/audit"
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
		nil,
		nil,
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

func TestLoginWritesAuditLog(t *testing.T) {
	repo := user.NewInMemoryRepository()
	userSvc := user.NewService(repo)
	if _, err := userSvc.CreateUser(context.Background(), user.CreateUserRequest{
		Username: "audited-user",
		Password: "secret123",
		Name:     "Audited User",
		Email:    "audited@example.com",
	}); err != nil {
		t.Fatalf("create user: %v", err)
	}

	auditSvc := audit.NewService(audit.NewInMemoryRepository(), "未知")
	svc := NewService(
		userSvc,
		jwt.NewManager("test-secret", 15*time.Minute, 24*time.Hour),
		nil,
		captchaStub{ok: true},
		auditSvc,
		nil,
	)

	ctx := audit.WithLoginMeta(context.Background(), audit.LoginMeta{
		RequestID: "req-1",
		IP:        "127.0.0.1",
		UserAgent: "unit-test",
	})
	if _, err := svc.Login(ctx, LoginRequest{
		Username:    "audited-user",
		Password:    "secret123",
		CaptchaID:   "captcha-id",
		CaptchaCode: "123456",
	}); err != nil {
		t.Fatalf("login: %v", err)
	}

	logs, err := auditSvc.ListLoginLogs(context.Background(), audit.ListLoginLogsQuery{})
	if err != nil {
		t.Fatalf("list login logs: %v", err)
	}
	if logs.Total != 1 {
		t.Fatalf("expected 1 login log, got %d", logs.Total)
	}
	if !logs.Items[0].Success || logs.Items[0].Username != "audited-user" {
		t.Fatalf("unexpected login log: %#v", logs.Items[0])
	}
}
