package main

import (
	"log"
	"net/http"

	"go-server/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server startup failed: %v", err)
	}
}
