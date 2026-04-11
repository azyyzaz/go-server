package config

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName string
	Env     string
	HTTP    HTTPConfig
	DB      DBConfig
	Redis   RedisConfig
	Log     LogConfig
	Audit   AuditConfig
	JWT     JWTConfig
}

type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
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

type AuditConfig struct {
	Enabled         bool
	RetentionDays   int
	CleanupInterval time.Duration
	RegionFallback  string
}

func (c HTTPConfig) Addr() string {
	return c.Host + ":" + c.Port
}

func Load() Config {
	return Config{
		AppName: viper.GetString("app.name"),
		Env:     viper.GetString("app.env"),
		HTTP: HTTPConfig{
			Host:         viper.GetString("http.host"),
			Port:         viper.GetString("http.port"),
			ReadTimeout:  viper.GetDuration("http.readTimeout"),
			WriteTimeout: viper.GetDuration("http.writeTimeout"),
		},
		DB: DBConfig{
			Enabled:         viper.GetBool("db.enabled"),
			DSN:             viper.GetString("db.dsn"),
			AutoMigrate:     viper.GetBool("db.autoMigrate"),
			MaxIdleConns:    viper.GetInt("db.maxIdleConns"),
			MaxOpenConns:    viper.GetInt("db.maxOpenConns"),
			ConnMaxLifetime: viper.GetDuration("db.connMaxLifetime"),
			ConnMaxIdleTime: viper.GetDuration("db.connMaxIdleTime"),
		},
		Redis: RedisConfig{
			Enabled:  viper.GetBool("redis.enabled"),
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		},
		Log: LogConfig{
			Level:      viper.GetString("log.level"),
			Filename:   viper.GetString("log.filename"),
			MaxSizeMB:  viper.GetInt("log.maxSizeMB"),
			MaxBackups: viper.GetInt("log.maxBackups"),
			MaxAgeDays: viper.GetInt("log.maxAgeDays"),
			Compress:   viper.GetBool("log.compress"),
		},
		Audit: AuditConfig{
			Enabled:         viper.GetBool("audit.enabled"),
			RetentionDays:   viper.GetInt("audit.retentionDays"),
			CleanupInterval: viper.GetDuration("audit.cleanupInterval"),
			RegionFallback:  viper.GetString("audit.regionFallback"),
		},
		JWT: JWTConfig{
			Secret:          viper.GetString("jwt.secret"),
			AccessTokenTTL:  viper.GetDuration("jwt.accessTokenTTL"),
			RefreshTokenTTL: viper.GetDuration("jwt.refreshTokenTTL"),
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
