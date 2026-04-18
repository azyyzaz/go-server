package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-server/internal/config"

	"github.com/gin-gonic/gin"
)

func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	allowMethods := joinOrDefault(cfg.AllowMethods, []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	allowHeaders := joinOrDefault(cfg.AllowHeaders, []string{"Authorization", "Content-Type", "X-Request-ID", "X-Trace-ID"})
	exposeHeaders := joinOrDefault(cfg.ExposeHeaders, []string{"X-Request-ID", "X-Trace-ID"})

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		if origin != "" {
			if allowedOrigin, ok := matchOrigin(origin, cfg.AllowOrigins); ok {
				c.Header("Access-Control-Allow-Origin", allowedOrigin)
				c.Header("Vary", "Origin")
			}
			c.Header("Access-Control-Allow-Methods", allowMethods)
			c.Header("Access-Control-Allow-Headers", allowHeaders)
			c.Header("Access-Control-Expose-Headers", exposeHeaders)
			if cfg.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", strconv.FormatInt(int64(cfg.MaxAge/time.Second), 10))
			}
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func matchOrigin(origin string, allowOrigins []string) (string, bool) {
	if len(allowOrigins) == 0 {
		return "*", true
	}
	for _, item := range allowOrigins {
		item = strings.TrimSpace(item)
		if item == "*" || strings.EqualFold(item, origin) {
			return item, true
		}
	}
	return "", false
}

func joinOrDefault(values []string, fallback []string) string {
	if len(values) == 0 {
		values = fallback
	}
	return strings.Join(values, ", ")
}
