package logger

import (
	"go-server/internal/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var global *zap.Logger

func Init(cfg config.LogConfig, env string) error {
	level := new(zap.AtomicLevel)
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level.SetLevel(zap.InfoLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "ts"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)
	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	})

	consoleWriter := zapcore.AddSync(os.Stdout)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	if env != "production" {
		options = append(options, zap.Development())
	}

	global = zap.New(core, options...)

	return nil
}

func L() *zap.Logger {
	if global == nil {
		return zap.NewNop()
	}
	return global
}

func S() *zap.SugaredLogger {
	return L().Sugar()
}

func Sync() error {
	return L().Sync()
}
