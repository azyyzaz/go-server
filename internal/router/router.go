package router

import (
	"go-server/internal/config"
	"go-server/internal/middleware"
	"go-server/internal/modules/health"
	"go-server/internal/modules/user"

	"github.com/gin-gonic/gin"
)

func New(cfg config.Config) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	api := r.Group("/api/v1")

	healthHandler := health.NewHandler()
	healthHandler.Register(api.Group("/health"))

	userRepo := user.NewInMemoryRepository()
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.Register(api.Group("/users"))

	return r
}
