package middlewares

import (
	"fmt"
	"time"

	"github.com/go-auth-service/src/config"
	"github.com/gofiber/fiber/v2"
)

func RateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {

		ctx := c.Context()

		userId, _ := c.Locals("userId").(string)

		var identifier string
		if userId != "" {
			identifier = userId
		} else {
			identifier = c.IP()
		}

		key := fmt.Sprintf("rate:%s:%s", identifier, c.Path())

		client, err := config.NewRedisClient(ctx)
		if err != nil {
			return err
		}

		window := 10 * time.Second
		maxRequests := int64(1000)

		count, err := client.Incr(ctx, key).Result()
		if err != nil {
			return err
		}

		if count == 1 {
			client.Expire(ctx, key, window)
		}

		if count > maxRequests {
			return c.Status(fiber.StatusTooManyRequests).
				JSON(fiber.Map{"message": "Too many requests, try again later"})
		}

		return c.Next()
	}
}
