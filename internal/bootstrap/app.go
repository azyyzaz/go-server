package bootstrap

import (
	"log"
	"net/http"

	"go-server/internal/config"
	"go-server/internal/router"
)

func Run() error {
	cfg := config.Load()
	engine := router.New(cfg)

	srv := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	log.Printf("[%s] HTTP server running on %s (%s)", cfg.AppName, cfg.HTTP.Addr(), cfg.Env)
	return srv.ListenAndServe()
}
