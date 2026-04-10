package user

import (
	"bytes"
	"context"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestExportUsersGeneratesExcelWorkbook(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), CreateUserRequest{
		Username: "alice",
		Password: "secret123",
		Name:     "Alice",
		Email:    "alice@example.com",
		Phone:    "13800138000",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	data, err := svc.ExportUsers(context.Background(), ListUsersQuery{})
	if err != nil {
		t.Fatalf("export users: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty excel data")
	}

	file, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("open excel workbook: %v", err)
	}
	defer func() { _ = file.Close() }()

	username, err := file.GetCellValue("Users", "B2")
	if err != nil {
		t.Fatalf("read username cell: %v", err)
	}
	if username != "alice" {
		t.Fatalf("expected exported username alice, got %q", username)
	}

	status, err := file.GetCellValue("Users", "F2")
	if err != nil {
		t.Fatalf("read status cell: %v", err)
	}
	if status != "启用" {
		t.Fatalf("expected exported status 启用, got %q", status)
	}
}
