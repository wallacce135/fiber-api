package router

import (
	"api/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Get("/hello", handler.HelloHandler)
	api.Get("/bye", handler.ByeWorldHandler)

	auth := app.Group("/auth")
	auth.Post("/login", handler.LoginHandler)
	auth.Post("/register", handler.RegisterHandler)

}
