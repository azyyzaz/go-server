// Package main Go Server API
//
//	@title			Go Server API
//	@version		1.0
//	@description	企业级后台管理系统 API
//	@host			localhost:8080
//	@BasePath		/api/v1
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				格式: Bearer <token>
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
