package config

import (
	"os"
	"time"
)

type Config struct {
	AppName string
	Env     string
	HTTP    HTTPConfig
}

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c HTTPConfig) Addr() string {
	return c.Host + ":" + c.Port
}

func Load() Config {
	return Config{
		AppName: getEnv("APP_NAME", "go-server"),
		Env:     getEnv("APP_ENV", "local"),
		HTTP: HTTPConfig{
			Host:         getEnv("HTTP_HOST", "0.0.0.0"),
			Port:         getEnv("HTTP_PORT", "8080"),
			ReadTimeout:  mustDuration(getEnv("HTTP_READ_TIMEOUT", "5s"), 5*time.Second),
			WriteTimeout: mustDuration(getEnv("HTTP_WRITE_TIMEOUT", "10s"), 10*time.Second),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustDuration(raw string, fallback time.Duration) time.Duration {
	d, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return d
}
