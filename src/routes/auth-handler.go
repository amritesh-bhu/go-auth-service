package routes

import (
	"context"
	"time"

	"github.com/go-auth-service/src/domain"
	"github.com/go-auth-service/src/entities"
	"github.com/go-auth-service/src/helper"
	"github.com/go-auth-service/src/middlewares"
	"github.com/go-auth-service/src/services"
	"github.com/gofiber/fiber/v2"
)

func AuthHandler(router fiber.Router) {

	auth := router.Group("/auth")

	auth.Post("/register", middlewares.RateLimiter(), func(c *fiber.Ctx) error {

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

	auth.Post("/login", middlewares.RateLimiter(), func(c *fiber.Ctx) error {

		ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
		defer cancel()

		var cred domain.LoginRequest

		if err := c.BodyParser(&cred); err != nil {
			c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		user, err := domain.LoginUser(ctx, &cred)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		userId := user.ID.Hex()

		accessToken, sessionId, err := helper.GetTokens(ctx, userId)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "sessionId",
			Value:    sessionId,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Strict",
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		return c.Status(200).JSON(fiber.Map{"accessToken": accessToken})
	})

	auth.Post("/logout", middlewares.UserSession(), middlewares.RateLimiter(), func(c *fiber.Ctx) error {

		sessionId := c.Locals("sessionId").(string)
		userId := c.Locals("userId").(string)
		ctx := c.Locals("ctx").(context.Context)

		err := services.DeleteRefreshToken(ctx, sessionId, userId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "sessionId",
			Value:    "",
			Secure:   true,
			HTTPOnly: true,
			SameSite: "None",
			Path:     "/",
			Expires:  time.Now().Add(-time.Hour),
		})
		return c.Status(200).JSON(fiber.Map{"message": "Logged out successfully!"})
	})

	auth.Post("/multiple-device-logout", middlewares.UserSession(), middlewares.RateLimiter(), func(c *fiber.Ctx) error {

		userId := c.Locals("userId").(string)
		ctx := c.Locals("ctx").(context.Context)

		err := services.LogOutFromMultipleDevice(ctx, userId)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err})
		}
		c.Cookie(&fiber.Cookie{
			Name:     "sessionId",
			Value:    "",
			Secure:   true,
			HTTPOnly: true,
			SameSite: "None",
			Expires:  time.Now().Add(-time.Hour),
		})
		return c.Status(200).JSON(fiber.Map{"message": "Logged out successfully!"})
	})

	auth.Post("/refresh-tokens", middlewares.UserSession(), middlewares.RateLimiter(), func(c *fiber.Ctx) error {

		userId := c.Locals("userId").(string)
		sessionId := c.Locals("sessionId").(string)
		ctx := c.Locals("ctx").(context.Context)

		err := services.DeleteRefreshToken(ctx, sessionId, userId)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err})
		}

		accessToken, sessionId, err := helper.GetTokens(ctx, userId)
		if err != nil {
			c.Status(500).JSON(fiber.Map{"error": err})
		}

		c.Cookie(&fiber.Cookie{
			Name:     "sessionId",
			Value:    sessionId,
			SameSite: "Strict",
			HTTPOnly: true,
			Secure:   true,
			Path:     "/",
			Expires:  time.Now().Add(7 * 24 * time.Hour),
		})

		return c.Status(200).JSON(fiber.Map{"accessToken": accessToken})
	})
}
