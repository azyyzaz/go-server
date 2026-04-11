package audit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-server/internal/middleware"
	"go-server/internal/modules/audit"

	"github.com/gin-gonic/gin"
)

func newAuditTestRouter(svc audit.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	audit.NewHandler(svc).Register(r.Group("/system/audits"))
	return r
}

func TestListOperationLogsRejectsInvalidTimeRange(t *testing.T) {
	svc := audit.NewService(audit.NewInMemoryRepository(), "未知")
	router := newAuditTestRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/system/audits/operation-logs?start_time=invalid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestListLoginLogsReturnsSuccess(t *testing.T) {
	svc := audit.NewService(audit.NewInMemoryRepository(), "未知")
	if err := svc.RecordLogin(context.Background(), audit.LoginLogEntry{
		RequestID: "req-1",
		Username:  "alice",
		IP:        "127.0.0.1",
		Success:   true,
	}); err != nil {
		t.Fatalf("record login: %v", err)
	}

	router := newAuditTestRouter(svc)
	req := httptest.NewRequest(http.MethodGet, "/system/audits/login-logs", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
