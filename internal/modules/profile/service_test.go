package profile

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"testing"

	"go-server/internal/config"
	"go-server/internal/modules/audit"
	filemodule "go-server/internal/modules/file"
	"go-server/internal/modules/user"
)

func TestProfileServiceSupportsProfileLifecycle(t *testing.T) {
	userSvc := user.NewService(user.NewInMemoryRepository())
	auditSvc := audit.NewService(audit.NewInMemoryRepository(), "test-region")
	fileSvc := newProfileFileService(t)
	svc := NewService(userSvc, auditSvc, fileSvc)

	created, err := userSvc.CreateUser(context.Background(), user.CreateUserRequest{
		Username: "alice",
		Password: "secret123",
		Name:     "Alice",
		Email:    "alice@example.com",
		Phone:    "13800138000",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	profile, err := svc.GetProfile(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get profile: %v", err)
	}
	if profile.Username != "alice" {
		t.Fatalf("expected username alice, got %q", profile.Username)
	}

	updated, err := svc.UpdateProfile(context.Background(), created.ID, UpdateProfileRequest{
		Name:  "Alice Doe",
		Email: "alice.doe@example.com",
		Phone: "13900139000",
	})
	if err != nil {
		t.Fatalf("update profile: %v", err)
	}
	if updated.Name != "Alice Doe" || updated.Email != "alice.doe@example.com" || updated.Phone != "13900139000" {
		t.Fatalf("unexpected updated profile: %+v", updated)
	}

	if err := svc.ChangePassword(context.Background(), created.ID, ChangePasswordRequest{
		OldPassword: "secret123",
		NewPassword: "newpass456",
	}); err != nil {
		t.Fatalf("change password: %v", err)
	}

	if err := userSvc.ChangePassword(context.Background(), created.ID, "secret123", "another789"); err == nil {
		t.Fatalf("expected old password to be invalid after change")
	}
	if err := userSvc.ChangePassword(context.Background(), created.ID, "newpass456", "another789"); err != nil {
		t.Fatalf("expected new password to work, got %v", err)
	}
}

func TestProfileServiceUploadsAvatarAndListsOwnLoginLogs(t *testing.T) {
	userSvc := user.NewService(user.NewInMemoryRepository())
	auditSvc := audit.NewService(audit.NewInMemoryRepository(), "test-region")
	fileSvc := newProfileFileService(t)
	svc := NewService(userSvc, auditSvc, fileSvc)

	created, err := userSvc.CreateUser(context.Background(), user.CreateUserRequest{
		Username: "bob",
		Password: "secret123",
		Name:     "Bob",
		Email:    "bob@example.com",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	fileHeader := makeProfileFileHeader(t, "avatar.png", []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	})
	uploaded, err := svc.UploadAvatar(context.Background(), created.ID, fileHeader)
	if err != nil {
		t.Fatalf("upload avatar: %v", err)
	}
	if uploaded.Avatar == "" {
		t.Fatalf("expected avatar url")
	}

	refreshed, err := userSvc.GetUser(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("get updated user: %v", err)
	}
	if refreshed.Avatar != uploaded.Avatar {
		t.Fatalf("expected avatar %q, got %q", uploaded.Avatar, refreshed.Avatar)
	}

	if err := auditSvc.RecordLogin(context.Background(), audit.LoginLogEntry{
		RequestID: "req-1",
		UserID:    &created.ID,
		Username:  "bob",
		IP:        "127.0.0.1",
		Success:   true,
	}); err != nil {
		t.Fatalf("record user login log: %v", err)
	}
	if err := auditSvc.RecordLogin(context.Background(), audit.LoginLogEntry{
		RequestID: "req-2",
		Username:  "other",
		IP:        "10.0.0.2",
		Success:   true,
	}); err != nil {
		t.Fatalf("record other login log: %v", err)
	}

	logs, err := svc.ListMyLoginLogs(context.Background(), created.ID, audit.ListLoginLogsQuery{})
	if err != nil {
		t.Fatalf("list my login logs: %v", err)
	}
	if logs.Total != 1 || len(logs.Items) != 1 {
		t.Fatalf("expected one own login log, got %+v", logs)
	}
	if logs.Items[0].Username != "bob" {
		t.Fatalf("expected bob login log, got %q", logs.Items[0].Username)
	}
}

func newProfileFileService(t *testing.T) filemodule.Service {
	t.Helper()

	storage := newProfileStorage(t)
	return filemodule.NewService(
		filemodule.NewInMemoryRepository(),
		storage,
		filemodule.NewValidator(10, []string{".png", ".jpg"}),
		filemodule.NewValidator(2, []string{".png", ".jpg"}),
	)
}

func newProfileStorage(t *testing.T) filemodule.Storage {
	t.Helper()

	return mustNewProfileStorage(config.FileConfig{
		Storage: "local",
		Local: config.LocalFileConfig{
			BaseDir:    t.TempDir(),
			BaseURL:    "/uploads",
			DateLayout: "2006/01/02",
		},
	})
}

func mustNewProfileStorage(cfg config.FileConfig) filemodule.Storage {
	storage, err := filemodule.NewStorage(cfg)
	if err != nil {
		panic(err)
	}
	return storage
}

func makeProfileFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	req := httptest.NewRequest("POST", "/profile/avatar", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(int64(len(body.Bytes()) + 1024)); err != nil {
		t.Fatalf("parse multipart form: %v", err)
	}

	fileHeader := req.MultipartForm.File["file"][0]
	fileHeader.Size = int64(len(content))
	return fileHeader
}
