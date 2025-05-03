package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	return cors.New(
		cors.Config{
			AllowOrigins: os.Getenv("CORS_ORIGIN"),
			AllowMethods: os.Getenv("CORS_METHODS"),
		},
	)
}
