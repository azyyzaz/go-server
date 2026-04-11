package dept

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func newDeptTestRouter(svc Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	handler := NewHandler(svc)
	handler.Register(r.Group("/system/depts"))
	return r
}

func TestCreateDeptRejectsInvalidEmail(t *testing.T) {
	repo := NewInMemoryRepository()
	svc := NewService(repo)
	router := newDeptTestRouter(svc)

	body, _ := json.Marshal(map[string]any{"name": "研发部", "email": "invalid", "status": 1})
	req := httptest.NewRequest(http.MethodPost, "/system/depts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
