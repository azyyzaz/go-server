package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func newUserTestRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	handler := NewHandler(svc)
	handler.Register(r.Group("/system/users"))
	return r
}

func TestUpdateUserStatusRejectsInvalidStatus(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	created, err := svc.CreateUser(context.Background(), CreateUserRequest{
		Username: "bob",
		Password: "secret123",
		Name:     "Bob",
		Email:    "bob@example.com",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	router := newUserTestRouter(svc)
	body, _ := json.Marshal(map[string]int{"status": 2})
	req := httptest.NewRequest(http.MethodPut, "/system/users/"+strconv.FormatUint(uint64(created.ID), 10)+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestExportUsersReturnsXLSXHeaders(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)

	_, err := svc.CreateUser(context.Background(), CreateUserRequest{
		Username: "carol",
		Password: "secret123",
		Name:     "Carol",
		Email:    "carol@example.com",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}

	router := newUserTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/system/users/export", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Disposition"); got != "attachment; filename=\"users.xlsx\"" {
		t.Fatalf("unexpected content disposition: %s", got)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		t.Fatalf("unexpected content type: %s", got)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty export body")
	}
}
