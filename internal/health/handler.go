package health

import (
	"database/sql"

	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func HealthCheck(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "Service is healthy",
		})
	}
}

func PanicTest(c *fiber.Ctx) error {
	panic("Panic test")
}

func ErrorLogTest(c *fiber.Ctx) error {
	err := sql.ErrNoRows
	log.Error("This is an error log test : ", err)
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"status":  "error",
		"message": "This is an error log test",
	})
}

func InfoLogTest(c *fiber.Ctx) error {
	log.Info("This is an info log test")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "This is an info log test",
	})
}

func DebugLogTest(c *fiber.Ctx) error {
	log.Debug("This is a debug log test")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "ok",
		"message": "This is a debug log test",
	})
}
