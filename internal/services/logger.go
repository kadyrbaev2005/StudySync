package services

import (
	"io"
	"log/slog"
	"os"

	"github.com/kadyrbayev2005/studysync/internal/utils"
)

// Logger — глобальный экземпляр логгера
var Logger *slog.Logger

// InitLogger инициализирует логгер по переменным окружения
// LOG_LEVEL: debug, info, warn, error (default: info)
// LOG_FORMAT: json, text (default: json)
// LOG_OUTPUT: stdout, stderr (default: stdout)
func InitLogger() {
	logLevel := utils.GetEnv("LOG_LEVEL", "info")
	logFormat := utils.GetEnv("LOG_FORMAT", "json")
	logOutput := utils.GetEnv("LOG_OUTPUT", "stdout")

	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{Level: level}

	var writer io.Writer
	switch logOutput {
	case "stderr":
		writer = os.Stderr
	default:
		writer = os.Stdout
	}

	var handler slog.Handler
	switch logFormat {
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}

// Удобные функции для логирования
func Debug(msg string, args ...any) {
	if Logger != nil {
		Logger.Debug(msg, args...)
	}
}

func Info(msg string, args ...any) {
	if Logger != nil {
		Logger.Info(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if Logger != nil {
		Logger.Warn(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if Logger != nil {
		Logger.Error(msg, args...)
	}
}
