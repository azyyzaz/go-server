package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName string
	Env     string
	HTTP    HTTPConfig
	DB      DBConfig
	Redis   RedisConfig
	Log     LogConfig
}

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	Enabled         bool
	DSN             string
	AutoMigrate     bool
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	MaxOpenConns    int
}

type RedisConfig struct {
	Enabled  bool
	Addr     string
	Password string
	DB       int
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
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
		DB: DBConfig{
			Enabled:         mustBool(getEnv("DB_ENABLED", "false"), false),
			DSN:             getEnv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/go_server?charset=utf8mb4&parseTime=True&loc=Local"),
			AutoMigrate:     mustBool(getEnv("DB_AUTO_MIGRATE", "true"), true),
			MaxIdleConns:    mustInt(getEnv("DB_MAX_IDLE_CONNS", "10"), 10),
			MaxOpenConns:    mustInt(getEnv("DB_MAX_OPEN_CONNS", "100"), 100),
			ConnMaxLifetime: mustDuration(getEnv("DB_CONN_MAX_LIFETIME", "1h"), time.Hour),
			ConnMaxIdleTime: mustDuration(getEnv("DB_CONN_MAX_IDLE_TIME", "30m"), 30*time.Minute),
		},
		Redis: RedisConfig{
			Enabled:  mustBool(getEnv("REDIS_ENABLED", "false"), false),
			Addr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       mustInt(getEnv("REDIS_DB", "0"), 0),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Filename:   getEnv("LOG_FILENAME", "logs/app.log"),
			MaxSizeMB:  mustInt(getEnv("LOG_MAX_SIZE_MB", "20"), 20),
			MaxBackups: mustInt(getEnv("LOG_MAX_BACKUPS", "10"), 10),
			MaxAgeDays: mustInt(getEnv("LOG_MAX_AGE_DAYS", "30"), 30),
			Compress:   mustBool(getEnv("LOG_COMPRESS", "false"), false),
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

func mustBool(raw string, fallback bool) bool {
	d, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return d
}

func mustInt(raw string, fallback int) int {
	d, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return d
}
