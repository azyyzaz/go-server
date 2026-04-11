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

// Ping godoc
//
//	@Summary		服务健康检查
//	@Description	用于确认服务是否正常启动并可对外提供请求处理能力。
//	@Tags			健康检查
//	@Produce		json
//	@Success		200	{object}	response.Body{data=map[string]interface{}}
//	@Router			/health [get]
func (h *Handler) Ping(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "up",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
