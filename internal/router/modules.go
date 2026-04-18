package router

import (
	"context"

	"go-server/internal/modules/audit"
	"go-server/internal/modules/auth"
	"go-server/internal/modules/captcha"
	"go-server/internal/modules/dashboard"
	"go-server/internal/modules/dept"
	"go-server/internal/modules/dict"
	filemodule "go-server/internal/modules/file"
	"go-server/internal/modules/health"
	"go-server/internal/modules/menu"
	"go-server/internal/modules/profile"
	"go-server/internal/modules/role"
	"go-server/internal/modules/user"

	"github.com/gin-gonic/gin"
)

type appServices struct {
	user      user.Service
	role      role.Service
	menu      menu.Service
	dept      dept.Service
	dict      dict.Service
	dashboard dashboard.Service
	file      filemodule.Service
}

func registerBaseRoutes(_ *gin.Engine, api *gin.RouterGroup, deps ModuleDeps, services appServices) {
	health.NewHandler(deps.DB, deps.Redis).Register(api.Group("/health"))

	captchaSvc := captcha.NewService()
	captcha.NewHandler(captchaSvc).Register(api.Group("/auth"))

	sessionTracker := dashboard.NewSessionTracker(deps.Redis)
	authService := auth.NewService(
		services.user,
		deps.JWTManager,
		deps.JWTBlacklist,
		captchaSvc,
		deps.AuditService,
		sessionTracker,
	)
	auth.NewHandler(authService).Register(api.Group("/auth"))
}

func registerSystemModules(api *gin.RouterGroup, deps ModuleDeps, services appServices) {
	if deps.AuditService != nil {
		audit.NewHandler(deps.AuditService).Register(api.Group("/system/audits"))
	}

	user.NewHandler(services.user).Register(api.Group("/system/users"))
	role.NewHandler(services.role).Register(api.Group("/system/roles"))
	menu.NewHandler(services.menu).Register(api.Group("/system/menus"))
	dept.NewHandler(services.dept).Register(api.Group("/system/depts"))
	dict.NewHandler(services.dict).Register(api.Group("/system/dicts"))
	filemodule.NewHandler(services.file).Register(api.Group("/system/files"))
}

func registerBusinessModules(api *gin.RouterGroup, deps ModuleDeps, services appServices) {
	dashboard.NewHandler(services.dashboard).Register(api.Group("/dashboard"))
	profileSvc := profile.NewService(services.user, deps.AuditService, services.file)
	profile.NewHandler(profileSvc).Register(api.Group("/profile"))
}

func buildServices(deps ModuleDeps) appServices {
	return appServices{
		user:      buildUserService(deps),
		role:      buildRoleService(deps),
		menu:      buildMenuService(deps),
		dept:      buildDeptService(deps),
		dict:      buildDictService(deps),
		dashboard: buildDashboardService(deps),
		file:      buildFileService(deps),
	}
}

func buildUserService(deps ModuleDeps) user.Service {
	var repo user.Repository
	if deps.DB != nil {
		repo = user.NewGORMRepository(deps.DB)
	} else {
		repo = user.NewInMemoryRepository()
	}

	svc := user.NewService(repo)
	if deps.DB == nil {
		_, _ = svc.CreateUser(context.Background(), user.CreateUserRequest{
			Username: "admin",
			Password: "admin123",
			Name:     "Admin",
			Email:    "admin@example.com",
		})
	}
	return svc
}

func buildRoleService(deps ModuleDeps) role.Service {
	var repo role.Repository
	if deps.DB != nil {
		repo = role.NewGORMRepository(deps.DB)
	} else {
		repo = role.NewInMemoryRepository()
	}
	return role.NewService(repo)
}

func buildMenuService(deps ModuleDeps) menu.Service {
	var repo menu.Repository
	if deps.DB != nil {
		repo = menu.NewGORMRepository(deps.DB)
	} else {
		repo = menu.NewInMemoryRepository()
	}
	return menu.NewService(repo)
}

func buildDeptService(deps ModuleDeps) dept.Service {
	var repo dept.Repository
	if deps.DB != nil {
		repo = dept.NewGORMRepository(deps.DB)
	} else {
		repo = dept.NewInMemoryRepository()
	}
	return dept.NewService(repo)
}

func buildDictService(deps ModuleDeps) dict.Service {
	var repo dict.Repository
	if deps.DB != nil {
		repo = dict.NewGORMRepository(deps.DB)
	} else {
		repo = dict.NewInMemoryRepository()
	}
	return dict.NewService(repo, deps.Redis)
}

func buildDashboardService(deps ModuleDeps) dashboard.Service {
	var repo dashboard.Repository
	if deps.DB != nil {
		repo = dashboard.NewGORMRepository(deps.DB, deps.Redis)
	} else {
		repo = dashboard.NewInMemoryRepository()
	}
	return dashboard.NewService(repo)
}

func buildFileService(deps ModuleDeps) filemodule.Service {
	var repo filemodule.Repository
	if deps.DB != nil {
		repo = filemodule.NewGORMRepository(deps.DB)
	} else {
		repo = filemodule.NewInMemoryRepository()
	}

	return filemodule.NewService(
		repo,
		deps.FileStorage,
		filemodule.NewValidator(deps.Config.File.MaxSizeMB, deps.Config.File.AllowedExts),
		filemodule.NewValidator(deps.Config.File.AvatarMaxSize, []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}),
	)
}
