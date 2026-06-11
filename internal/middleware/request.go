package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set("X-Request-Id", requestID)
		c.Locals("requestId", requestID)
		return c.Next()
	}
}

func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		logger.Info("http request",
			zap.String("request_id", localString(c, "requestId")),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
		)

		return err
	}
}

func localString(c *fiber.Ctx, key string) string {
	value, ok := c.Locals(key).(string)
	if !ok {
		return ""
	}
	return value
}
