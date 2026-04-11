package dict

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func newDictTestRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	handler := NewHandler(svc)
	handler.Register(r.Group("/system/dicts"))
	return r
}

func TestCreateTypeRejectsInvalidStatus(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo, nil)
	router := newDictTestRouter(svc)

	body, _ := json.Marshal(map[string]any{"name": "状态", "code": "sys_status", "status": 2})
	req := httptest.NewRequest(http.MethodPost, "/system/dicts/types", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
