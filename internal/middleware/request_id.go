package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyRequestID = "request_id"
	ContextKeyTraceID   = "trace_id"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := strings.TrimSpace(c.GetHeader("X-Trace-ID"))
		if rid == "" {
			rid = strings.TrimSpace(c.GetHeader("X-Request-ID"))
		}
		if rid == "" {
			rid = generateRequestID()
		}
		c.Set(ContextKeyRequestID, rid)
		c.Set(ContextKeyTraceID, rid)
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Writer.Header().Set("X-Trace-ID", rid)
		c.Next()
	}
}

func generateRequestID() string {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return time.Now().UTC().Format("20060102150405.000000")
	}
	return time.Now().UTC().Format("20060102150405") + "-" + hex.EncodeToString(buf)
}
