package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-server/internal/modules/audit"

	"github.com/gin-gonic/gin"
)

func newAuditMiddlewareRouter(svc audit.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(AuditLog(svc))
	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func TestAuditMiddlewareRecordsWhitelistedPath(t *testing.T) {
	svc := audit.NewService(audit.NewInMemoryRepository(), "未知")
	router := newAuditMiddlewareRouter(svc)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"alice","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	result, err := svc.ListOperationLogs(req.Context(), audit.ListOperationLogsQuery{})
	if err != nil {
		t.Fatalf("list operation logs: %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("expected 1 operation log, got %d", result.Total)
	}
	if strings.Contains(result.Items[0].RequestBody, "secret") {
		t.Fatalf("expected password to be redacted, got %s", result.Items[0].RequestBody)
	}
}

func TestAuditMiddlewareSkipsNonWhitelistedPath(t *testing.T) {
	svc := audit.NewService(audit.NewInMemoryRepository(), "未知")
	router := newAuditMiddlewareRouter(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	result, err := svc.ListOperationLogs(req.Context(), audit.ListOperationLogsQuery{})
	if err != nil {
		t.Fatalf("list operation logs: %v", err)
	}
	if result.Total != 0 {
		t.Fatalf("expected 0 operation logs, got %d", result.Total)
	}
}
