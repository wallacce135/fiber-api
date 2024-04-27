package middleware

import (
	"api/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(context *fiber.Ctx, err error) error {
	if err.Error() == "Mission or malformed JWT" {
		return context.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return context.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})
}
