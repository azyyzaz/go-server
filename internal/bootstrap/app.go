package bootstrap

import (
	"net/http"

	"go-server/internal/config"
	"go-server/internal/db"
	"go-server/internal/logger"
	"go-server/internal/redis"
	"go-server/internal/router"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func Run() error {
	godotenv.Load()
	cfg := config.Load()

	if err := logger.Init(cfg.Log, cfg.Env); err != nil {
		return err
	}

	defer logger.Sync()

	if cfg.DB.Enabled {
		_, err := db.Init(cfg.DB)
		if err != nil {
			// log.Fatal(err)
			logger.L().Fatal("database init failed", zap.Error(err))
		}
	}

	if cfg.Redis.Enabled {
		_, err := redis.Init(cfg.Redis)
		if err != nil {
			// log.Fatal(err)
			logger.L().Fatal("Redis init failed", zap.Error(err))
		}
	}

	engine := router.New(cfg)

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	// log.Printf("[%s] HTTP server running on %s (%s)", cfg.AppName, cfg.HTTP.Addr(), cfg.Env)
	logger.L().Info("http server running",
		zap.String("app", cfg.AppName),
		zap.String("Addr", cfg.HTTP.Addr()),
		zap.String("env", cfg.Env),
	)
	return srv.ListenAndServe()
}
