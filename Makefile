APP_NAME := go-server
BUILD_DIR := bin
BINARY := $(BUILD_DIR)/$(APP_NAME)
DB_URL ?= mysql://root:root@tcp(127.0.0.1:3306)/go_server

.PHONY: run build test migrate

run:
	go run ./cmd/api

build:
	go build -o $(BINARY) ./cmd/api

test:
	go test ./...

migrate:
	migrate -path ./migrations -database "$(DB_URL)" up
