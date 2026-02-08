package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func StructuredLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Início do cronômetro
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Processa a requisição
		c.Next()

		// Fim do cronômetro
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Monta os atributos do log
		attributes := []any{
			slog.String("method", method),
			slog.String("path", path),
			slog.Int("status", statusCode),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
		}

		// Adiciona User ID se estiver disponível (ex: depois do login)
		if userID, exists := c.Get("user_id"); exists {
			attributes = append(attributes, slog.Any("user_id", userID))
		}

		// Loga com nível apropriado
		if statusCode >= 500 {
			slog.Error("request_failed", attributes...)
		} else if statusCode >= 400 {
			slog.Warn("client_error", attributes...)
		} else {
			slog.Info("request_success", attributes...)
		}
	}
}