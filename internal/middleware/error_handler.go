package middleware

import (
	"errors"

	"go-server/internal/errcode"
	"go-server/internal/response"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Written() {
			return
		}
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var appErr *response.AppError
		if errors.As(err, &appErr) {
			response.Fail(c, appErr.Status, appErr.Code, appErr.Message)
			return
		}

		fallback := errcode.ErrInternalError.AsError()
		response.Fail(c, fallback.Status, fallback.Code, fallback.Message)
	}
}
