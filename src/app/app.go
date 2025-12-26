package app

import (
	"github.com/gofiber/fiber/v2"
)

func NewServer() *fiber.App {
	return fiber.New()
}
