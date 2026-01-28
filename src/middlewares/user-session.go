package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/go-auth-service/src/config"
	"github.com/go-auth-service/src/helper"
	"github.com/gofiber/fiber/v2"
)

func UserSession() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		c.Locals("ctx", ctx)

		sessionId := c.Cookies("sessionId")
		if sessionId == "" {
			c.Status(500).JSON(fiber.Map{"error": "Session Expired, please login again!"})
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := helper.ValidateKey(token, config.Load().AccessSecret)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token!"})
		}

		c.Locals("sessionId", sessionId)
		c.Locals("userId", claims.UserId)

		return c.Next()
	}
}
