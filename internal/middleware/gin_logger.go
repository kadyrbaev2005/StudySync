package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kadyrbayev2005/studysync/internal/services"
)

// GinLogger создает middleware для логирования HTTP-запросов с использованием slog
func GinLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Логируем через slog
		duration := param.Latency
		clientIP := param.ClientIP
		method := param.Method
		path := param.Path
		statusCode := param.StatusCode
		errorMessage := param.ErrorMessage

		if errorMessage != "" {
			services.Error("HTTP request error",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", duration,
				"client_ip", clientIP,
				"error", errorMessage,
			)
		} else {
			services.Info("HTTP request",
				"method", method,
				"path", path,
				"status", statusCode,
				"latency", duration.String(),
				"client_ip", clientIP,
			)
		}

		// Возвращаем пустую строку, так как логируем сами
		return ""
	})
}
