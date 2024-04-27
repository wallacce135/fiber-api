package main

import (
	"api/database"
	router "api/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "App Name",
	})

	router.SetupRouter(app)

	database.ConnectDB()

	log.Fatal(app.Listen(":4000"))
}
