package bootstrap

import (
	"net/http"
	"strings"

	"go-server/internal/config"
	"go-server/internal/db"
	"go-server/internal/logger"
	"go-server/internal/redis"
	"go-server/internal/router"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Run() error {

	// 初始化 Viper
	viper.SetConfigName("config") // 文件名（不含扩展名）
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // 在项目根目录查找

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // app.env → APP_ENV
	viper.AutomaticEnv()                                   // 环境变量自动覆盖 yaml 配置

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.L().Info("config reloaded", zap.String("file", e.Name))
	})

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
