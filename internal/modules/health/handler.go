package health

import (
	"context"
	"net/http"
	"time"

	"go-server/internal/response"

	"github.com/gin-gonic/gin"
	rdb "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	dependencyHealthy = "up"
	dependencyFailed  = "down"
	dependencySkipped = "disabled"
)

type Handler struct {
	db    *gorm.DB
	redis *rdb.Client
}

func NewHandler(db *gorm.DB, redis *rdb.Client) *Handler {
	return &Handler{db: db, redis: redis}
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("", h.Ping)
	rg.GET("/ready", h.Ready)
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

func (h *Handler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	checks := gin.H{
		"database": h.databaseStatus(ctx),
		"redis":    h.redisStatus(ctx),
	}

	statusCode := http.StatusOK
	status := dependencyHealthy
	for _, item := range checks {
		if item == dependencyFailed {
			statusCode = http.StatusServiceUnavailable
			status = dependencyFailed
			break
		}
	}

	c.JSON(statusCode, response.Body{
		Code:    response.CodeOK,
		Message: "success",
		Data: gin.H{
			"status":       status,
			"timestamp":    time.Now().UTC().Format(time.RFC3339),
			"dependencies": checks,
		},
		RequestID: response.RequestIDFromContext(c),
		TraceID:   response.TraceIDFromContext(c),
	})
}

func (h *Handler) databaseStatus(ctx context.Context) string {
	if h.db == nil {
		return dependencySkipped
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return dependencyFailed
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		return dependencyFailed
	}
	return dependencyHealthy
}

func (h *Handler) redisStatus(ctx context.Context) string {
	if h.redis == nil {
		return dependencySkipped
	}

	if err := h.redis.Ping(ctx).Err(); err != nil {
		return dependencyFailed
	}
	return dependencyHealthy
}
