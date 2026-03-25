package main

import (
	"net/http"

	"go-server/internal/bootstrap"
	"go-server/internal/logger"

	"go.uber.org/zap"
)

func main() {
	if err := bootstrap.Run(); err != nil && err != http.ErrServerClosed {
		// log.Fatalf("server startup failed: %v", err)
		logger.L().Fatal("server startup failed", zap.Error(err))
	}
}
