package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-Id"

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try to use an existing request header, otherwise generate a new one.
		rid := c.Get(RequestIDHeader)
		if rid == "" {
			rid = uuid.New().String()
		}

		// Make the request id available to later middleware/handlers via Locals
		c.Locals(RequestIDHeader, rid)

		// Expose it on the response headers as well
		c.Set(RequestIDHeader, rid)

		return c.Next()
	}
}
