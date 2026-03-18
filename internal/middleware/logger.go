package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		rid, _ := c.Get("request_id")
		latency := time.Since(start)
		log.Printf("rid=%v status=%d method=%s path=%s latency=%s ip=%s",
			rid,
			c.Writer.Status(),
			c.Request.Method,
			c.Request.URL.Path,
			latency,
			c.ClientIP(),
		)
	}
}
