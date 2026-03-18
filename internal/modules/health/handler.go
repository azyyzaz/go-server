package health

import (
	"time"

	"go-server/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.Ping)
}

func (h *Handler) Ping(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "up",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
