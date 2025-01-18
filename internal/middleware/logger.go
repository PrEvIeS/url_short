package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Начало времени выполнения запроса
		start := time.Now()

		// Продолжаем выполнение следующих middleware и обработчиков
		c.Next()

		// Вычисляем время выполнения запроса
		latency := time.Since(start)

		// Логируем сведения о запросе и ответе
		logger.Info("Request details",
			zap.String("URI", c.Request.RequestURI),
			zap.String("Method", c.Request.Method),
			zap.Duration("Latency", latency),
			zap.Int("Status", c.Writer.Status()),
			zap.Int("Response Size", c.Writer.Size()),
		)
	}
}
