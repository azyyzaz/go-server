// Package main Go Server 后台管理 API
//
//	@title			Go Server 后台管理 API
//	@version		1.0
//	@description	企业级后台管理系统接口文档，覆盖认证、用户、角色、菜单、部门、字典和审计等业务模块。
//	@host			localhost:8080
//	@BasePath		/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				请填写 Bearer Token，格式为：Bearer <token>
package main

import (
	"net/http"

	_ "go-server/docs"
	"go-server/internal/bootstrap"
	"go-server/internal/logger"

	"go.uber.org/zap"
)

func main() {
	if err := bootstrap.Run(); err != nil && err != http.ErrServerClosed {
		logger.L().Fatal("server startup failed", zap.Error(err))
	}
}
