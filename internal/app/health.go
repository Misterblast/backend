package app

import (
	"database/sql"

	"github.com/ghulammuzz/misterblast/internal/health"
	"github.com/gofiber/fiber/v2"
)

func RegisterHealthRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/hc", health.HealthCheck(db))
	app.Get("/panic", health.PanicTest)
	app.Get("/error", health.ErrorLogTest)
	app.Get("/info", health.InfoLogTest)
	app.Get("/debug", health.DebugLogTest)
}
