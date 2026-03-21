package bootstrap

import (
	"log"
	"net/http"

	"go-server/internal/config"
	"go-server/internal/db"
	"go-server/internal/router"

	"github.com/joho/godotenv"
)

func Run() error {
	godotenv.Load()
	cfg := config.Load()

	if cfg.DB.Enabled {
		_, err := db.Init(cfg.DB)
		if err != nil {
			log.Fatal(err)
		}
	}
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
