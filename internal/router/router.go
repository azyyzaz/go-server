package router

import (
	"context"

	"go-server/internal/config"
	appjwt "go-server/internal/jwt"
	"go-server/internal/middleware"
	"go-server/internal/modules/auth"
	"go-server/internal/modules/captcha"
	"go-server/internal/modules/dept"
	"go-server/internal/modules/dict"
	"go-server/internal/modules/health"
	"go-server/internal/modules/menu"
	"go-server/internal/modules/role"
	"go-server/internal/modules/user"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	rdb "github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func New(cfg config.Config, db *gorm.DB, redisClient *rdb.Client, jwtManager *appjwt.Manager, jwtBlacklist *appjwt.Blacklist, casbinEnforcer *casbin.Enforcer) *gin.Engine {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())

	// Swagger UI（不走鉴权）
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	api.Use(middleware.JWTAuth(jwtManager, jwtBlacklist,
		"/api/v1/auth",
		"/api/v1/health",
	))
	api.Use(middleware.CasbinAuth(casbinEnforcer,
		"/api/v1/auth",
		"/api/v1/health",
	))

	healthHandler := health.NewHandler()
	healthHandler.Register(api.Group("/health"))

	// 选择 repository：有 DB 用 GORM，否则降级内存
	var userRepo user.Repository
	if db != nil {
		userRepo = user.NewGORMRepository(db)
	} else {
		userRepo = user.NewInMemoryRepository()
	}
	userService := user.NewService(userRepo)

	// 无 DB 时预置 admin，方便本地调试
	if db == nil {
		_, _ = userService.CreateUser(context.Background(), user.CreateUserRequest{
			Username: "admin",
			Password: "admin123",
			Name:     "Admin",
			Email:    "admin@example.com",
		})
	}

	userHandler := user.NewHandler(userService)
	userHandler.Register(api.Group("/system/users"))

	var roleRepo role.Repository
	if db != nil {
		roleRepo = role.NewGORMRepository(db)
	} else {
		roleRepo = role.NewInMemoryRepository()
	}
	roleService := role.NewService(roleRepo)
	roleHandler := role.NewHandler(roleService)
	roleHandler.Register(api.Group("/system/roles"))

	var menuRepo menu.Repository
	if db != nil {
		menuRepo = menu.NewGORMRepository(db)
	} else {
		menuRepo = menu.NewInMemoryRepository()
	}
	menuService := menu.NewService(menuRepo)
	menuHandler := menu.NewHandler(menuService)
	menuHandler.Register(api.Group("/system/menus"))

	var deptRepo dept.Repository
	if db != nil {
		deptRepo = dept.NewGORMRepository(db)
	} else {
		deptRepo = dept.NewInMemoryRepository()
	}
	deptService := dept.NewService(deptRepo)
	deptHandler := dept.NewHandler(deptService)
	deptHandler.Register(api.Group("/system/depts"))

	var dictRepo dict.Repository
	if db != nil {
		dictRepo = dict.NewGORMRepository(db)
	} else {
		dictRepo = dict.NewInMemoryRepository()
	}
	dictService := dict.NewService(dictRepo, redisClient)
	dictHandler := dict.NewHandler(dictService)
	dictHandler.Register(api.Group("/system/dicts"))

	captchaSvc := captcha.NewService()
	captchaHandler := captcha.NewHandler(captchaSvc)
	captchaHandler.Register(api.Group("/auth"))

	authService := auth.NewService(userService, jwtManager, jwtBlacklist, captchaSvc)
	authHandler := auth.NewHandler(authService)
	authHandler.Register(api.Group("/auth"))

	return r
}
