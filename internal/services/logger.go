package services

import (
    "io"
    "log/slog"
    "os"
    "path/filepath"

    "github.com/kadyrbayev2005/studysync/internal/utils"
    "gopkg.in/natefinch/lumberjack.v2"
)

// Logger — глобальный экземпляр логгера
var Logger *slog.Logger

// InitLogger инициализирует логгер по переменным окружения
func InitLogger() {
    logLevel := utils.GetEnv("LOG_LEVEL", "info")
    logFormat := utils.GetEnv("LOG_FORMAT", "json")
    logOutput := utils.GetEnv("LOG_OUTPUT", "file")

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

    var writers []io.Writer
    
    switch logOutput {
    case "both":
        writers = append(writers, os.Stdout)
        fallthrough
    case "file":
        logsDir := "logs"
        if err := os.MkdirAll(logsDir, 0755); err != nil {
            panic(err)
        }
        
        logFile := &lumberjack.Logger{
            Filename:   filepath.Join(logsDir, "studysync.log"),
            MaxSize:    10,    // megabytes
            MaxBackups: 5,
            MaxAge:     30,    // days
            Compress:   true,
        }
        writers = append(writers, logFile)
    default: // stdout
        writers = append(writers, os.Stdout)
    }

    multiWriter := io.MultiWriter(writers...)

    var handler slog.Handler
    switch logFormat {
    case "text":
        handler = slog.NewTextHandler(multiWriter, opts)
    default:
        handler = slog.NewJSONHandler(multiWriter, opts)
    }

    Logger = slog.New(handler)
    slog.SetDefault(Logger)
    
    Info("Logger initialized", 
        "level", logLevel, 
        "format", logFormat, 
        "output", logOutput)
}

// Удобные функции для логирования - они ДОЛЖНЫ быть экспортированы
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