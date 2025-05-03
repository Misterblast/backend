package app

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/hc", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(500).SendString("Database not healthy")
		}
		return c.SendString("OK")
	})

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})
}
