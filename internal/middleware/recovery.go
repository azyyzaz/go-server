package middleware

import (
	"go-server/internal/errcode"
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
		e := errcode.ErrInternalError.AsError()
		response.Fail(c, e.Status, e.Code, e.Message)
	})
}
