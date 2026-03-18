package middleware

import (
	"log"
	"net/http"

	"go-server/internal/response"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Printf("panic recovered: %v", recovered)
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError, "internal server error")
	})
}
