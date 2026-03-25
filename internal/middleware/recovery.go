package middleware

import (
	"net/http"

	"go-server/internal/logger"
	"go-server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		// log.Printf("panic recovered: %v", recovered)
		logger.L().Error("panic recovered",
			zap.Any("error", recovered),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
		)
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
	})
}
