package middleware

import (
	"go-server/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		rid := c.GetString("request_id")
		latency := time.Since(start)
		status := c.Writer.Status()
		// log.Printf("rid=%v status=%d method=%s path=%s latency=%s ip=%s",
		// 	rid,
		// 	c.Writer.Status(),
		// 	c.Request.Method,
		// 	c.Request.URL.Path,
		// 	latency,
		// 	c.ClientIP(),
		// )

		fields := []zap.Field{
			zap.String("request_id", rid),
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.String("ip", c.ClientIP()),
			zap.String("error", c.Errors.String()),
			zap.Duration("latency", latency),
		}

		switch {
		case status >= 500:
			logger.L().Error("http request", fields...)
		case status >= 400:
			logger.L().Warn("http request", fields...)
		default:
			logger.L().Info("http request", fields...)
		}

	}
}
