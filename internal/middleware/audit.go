package middleware

import (
	"bytes"
	"io"
	"strings"
	"time"

	"go-server/internal/logger"
	"go-server/internal/modules/audit"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuditLog(svc audit.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if svc == nil || !shouldAuditPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		var bodyCopy []byte
		if c.Request.Body != nil {
			bodyCopy, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
		}

		start := time.Now()
		c.Next()

		username := c.GetString(ContextKeyUsername)
		if username == "" {
			username = audit.ExtractUsername(bodyCopy)
		}

		var userID *uint
		if raw, ok := c.Get(ContextKeyUserID); ok {
			if parsed, ok := raw.(uint); ok {
				userID = &parsed
			}
		}

		err := svc.RecordOperation(c.Request.Context(), audit.OperationLogEntry{
			RequestID:    c.GetString("request_id"),
			UserID:       userID,
			Username:     username,
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			IP:           c.ClientIP(),
			StatusCode:   c.Writer.Status(),
			RequestBody:  audit.SummarizeRequestBody(bodyCopy),
			ErrorMessage: c.Errors.String(),
			LatencyMS:    time.Since(start).Milliseconds(),
		})
		if err != nil {
			logger.L().Warn("record audit operation log failed", zap.Error(err))
		}
	}
}

func shouldAuditPath(path string) bool {
	if strings.HasPrefix(path, "/api/v1/system/") {
		return true
	}
	switch path {
	case "/api/v1/auth/login", "/api/v1/auth/logout", "/api/v1/auth/refresh":
		return true
	default:
		return false
	}
}
