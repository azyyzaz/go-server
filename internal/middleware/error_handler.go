package middleware

import (
	"errors"
	"net/http"

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

		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
	}
}
