package handler

import (
	"github.com/gofiber/fiber/v2"
)

func HelloHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "Hello World!"})
}

func ByeWorldHandler(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"status": "Bye World!"})
}
