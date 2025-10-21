package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start)
		// prefer request id from locals (set by RequestID middleware), fall back to header
		rid := ""
		if v := c.Locals(RequestIDHeader); v != nil {
			if s, ok := v.(string); ok {
				rid = s
			}
		}
		if rid == "" {
			rid = c.Get(RequestIDHeader)
		}

		logger.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", elapsed),
			zap.String("request_id", rid),
		)
		return err
	}
}
