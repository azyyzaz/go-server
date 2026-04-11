package bootstrap

import (
	"net/http"
	"strings"

	appcasbin "go-server/internal/casbin"
	"go-server/internal/config"
	"go-server/internal/db"
	appjwt "go-server/internal/jwt"
	"go-server/internal/logger"
	appredis "go-server/internal/redis"
	"go-server/internal/router"

	"github.com/fsnotify/fsnotify"
	rdb "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

	var gormDB *gorm.DB
	if cfg.DB.Enabled {
		var err error
		gormDB, err = db.Init(cfg.DB)
		if err != nil {
			logger.L().Fatal("database init failed", zap.Error(err))
		}
	}

	var redisClient *rdb.Client
	if cfg.Redis.Enabled {
		var err error
		redisClient, err = appredis.Init(cfg.Redis)
		if err != nil {
			logger.L().Fatal("Redis init failed", zap.Error(err))
		}
	}

	jwtManager := appjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTokenTTL, cfg.JWT.RefreshTokenTTL)
	var jwtBlacklist *appjwt.Blacklist
	if redisClient != nil {
		jwtBlacklist = appjwt.NewBlacklist(redisClient)
	}

	casbinEnforcer, err := appcasbin.Init("configs/rbac_model.conf", "configs/rbac_policy.csv", gormDB)
	if err != nil {
		logger.L().Fatal("casbin init failed", zap.Error(err))
	}

	engine := router.New(cfg, gormDB, jwtManager, jwtBlacklist, casbinEnforcer)

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	logger.L().Info("http server running",
		zap.String("app", cfg.AppName),
		zap.String("Addr", cfg.HTTP.Addr()),
		zap.String("env", cfg.Env),
	)
	return srv.ListenAndServe()
}
