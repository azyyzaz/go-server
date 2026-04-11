package role

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	appcasbin "go-server/internal/casbin"
	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func newRoleTestRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	handler := NewHandler(svc)
	handler.Register(r.Group("/system/roles"))
	return r
}

func TestCreateRoleRejectsInvalidStatus(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)
	router := newRoleTestRouter(svc)

	body, _ := json.Marshal(map[string]any{"name": "Admin", "code": "admin", "status": 2})
	req := httptest.NewRequest(http.MethodPost, "/system/roles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestListRolesReturnsCreatedRole(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)
	if _, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Admin", Code: "admin", Status: 1}); err != nil {
		t.Fatalf("create role: %v", err)
	}

	router := newRoleTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/system/roles", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestUpdateRoleAPIsEndpoint(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)
	enforcer := newTestEnforcer(t)
	appcasbin.Set(enforcer)
	defer appcasbin.Set(nil)

	created, err := svc.CreateRole(context.Background(), CreateRoleRequest{Name: "Operator", Code: "operator", Status: 1})
	if err != nil {
		t.Fatalf("create role: %v", err)
	}

	router := newRoleTestRouter(svc)
	body, _ := json.Marshal(map[string]any{
		"permissions": []map[string]string{{"path": "/api/v1/system/roles", "method": "GET"}},
	})
	req := httptest.NewRequest(http.MethodPut, "/system/roles/"+strconv.FormatUint(uint64(created.ID), 10)+"/apis", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestListRoleUsersRejectsInvalidRoleID(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)
	router := newRoleTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/system/roles/abc/users", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
