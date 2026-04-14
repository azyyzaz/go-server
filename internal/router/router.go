package router

import (
	"go-server/internal/config"
	appjwt "go-server/internal/jwt"
	"go-server/internal/middleware"
	"go-server/internal/modules/audit"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	rdb "github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type ModuleDeps struct {
	Config         config.Config
	DB             *gorm.DB
	Redis          *rdb.Client
	JWTManager     *appjwt.Manager
	JWTBlacklist   *appjwt.Blacklist
	CasbinEnforcer *casbin.Enforcer
	AuditService   audit.Service
}

func New(cfg config.Config, db *gorm.DB, redisClient *rdb.Client, jwtManager *appjwt.Manager, jwtBlacklist *appjwt.Blacklist, casbinEnforcer *casbin.Enforcer, auditSvc audit.Service) *gin.Engine {
	deps := ModuleDeps{
		Config:         cfg,
		DB:             db,
		Redis:          redisClient,
		JWTManager:     jwtManager,
		JWTBlacklist:   jwtBlacklist,
		CasbinEnforcer: casbinEnforcer,
		AuditService:   auditSvc,
	}
	services := buildServices(deps)

	r := setupEngine(deps)
	registerDocsRoutes(r)

	api := r.Group("/api/v1")
	applyAPIMiddleware(api, deps)

	registerBaseRoutes(r, api, deps, services)
	registerSystemModules(api, deps, services)
	registerBusinessModules(api, deps, services)

	return r
}

func setupEngine(deps ModuleDeps) *gin.Engine {
	if deps.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.ErrorHandler())
	r.Use(middleware.AuditLog(deps.AuditService))
	return r
}

func registerDocsRoutes(r *gin.Engine) {
	r.Static("/uploads", "./uploads")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func applyAPIMiddleware(api *gin.RouterGroup, deps ModuleDeps) {
	api.Use(middleware.JWTAuth(deps.JWTManager, deps.JWTBlacklist,
		"/api/v1/auth",
		"/api/v1/health",
	))
	api.Use(middleware.CasbinAuth(deps.CasbinEnforcer,
		"/api/v1/auth",
		"/api/v1/health",
		"/api/v1/profile",
	))
}
