package routes

import (
	"context"
	"time"

	"github.com/go-auth-service/src/domain"
	"github.com/go-auth-service/src/entities"
	"github.com/gofiber/fiber/v2"
)

func AuthHandler(router fiber.Router) {

	auth := router.Group("/auth")

	auth.Post("/register", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		var user entities.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		result, err := domain.RegisterUser(ctx, &user)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(200).JSON(fiber.Map{"message": "User created successfully!", "data": result})
	})

	auth.Post("/login", func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		var cred domain.LoginRequest

		if err := c.BodyParser(&cred); err != nil {
			c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		msg, err := domain.LoginUser(ctx, &cred)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(200).JSON(fiber.Map{"message": msg})
	})
}
