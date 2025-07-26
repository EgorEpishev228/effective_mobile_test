package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestBody := string(c.Body())
		err := c.Next()
		duration := time.Since(start)

		logger.WithFields(logrus.Fields{
			"method":        c.Method(),
			"path":          c.Path(),
			"status":        c.Response().StatusCode(),
			"duration_ms":   duration.Milliseconds(),
			"ip":            c.IP(),
			"user_agent":    string(c.Request().Header.UserAgent()),
			"request_body":  requestBody,
			"response_body": string(c.Response().Body()),
			"content_type":  string(c.Response().Header.ContentType()),
		}).Info("HTTP Request")

		return err
	}
}
