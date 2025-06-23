package app

import (
	"database/sql"
	"fmt"
	"os"

	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"

	class "github.com/ghulammuzz/misterblast/internal/class/di"
	content "github.com/ghulammuzz/misterblast/internal/content/di"
	email "github.com/ghulammuzz/misterblast/internal/email/di"
	lesson "github.com/ghulammuzz/misterblast/internal/lesson/di"
	question "github.com/ghulammuzz/misterblast/internal/question/di"
	quiz "github.com/ghulammuzz/misterblast/internal/quiz/di"
	set "github.com/ghulammuzz/misterblast/internal/set/di"
	task "github.com/ghulammuzz/misterblast/internal/task/di"
	user "github.com/ghulammuzz/misterblast/internal/user/di"
)

func SetupRouter(db *sql.DB, redis *redis.Client) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Use(m.Cors())
	app.Use(m.Recover())
	app.Use(m.Metrics())

	api := app.Group("/v1")
	// api := app.Group("/v2")

	class.InitializedClassService(db, redis).Router(api)
	lesson.InitializedLessonService(db, redis, m.Validate).Router(api)
	set.InitializedSetService(db, redis, m.Validate).Router(api)
	question.InitializedQuestionService(db, redis, m.Validate).Router(api)
	user.InitializedUserService(db, m.Validate).Router(api)
	email.InitializedEmailService(db, m.Validate).Router(api)
	quiz.InitializedQuizService(db, m.Validate).Router(api)
	task.InitializeTaskService(db, m.Validate).Router(api)
	task.InitializeTaskSubmissionService(db, m.Validate).Router(api)
	content.InitializedContentService(db, redis, m.Validate).Router(api)
	content.InitializedAuthorService(db, redis, m.Validate).Router(api)

	app.Get("/.well-known/assetlinks.json", func(c *fiber.Ctx) error {
		jsonData, err := os.ReadFile("internal-link.json")
		if err != nil {
			log.Errorf("Failed to read assetlinks data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read assetlinks data",
			})
		}

		c.Set("Content-Type", "application/json")
		return c.Send(jsonData)
	})
	app.Get("/update-password", func(c *fiber.Ctx) error {
		return c.SendString(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>Open in App</title>
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
			</head>
			<body>
				<h3>Please open this link in the Go Assessment mobile app.</h3>
				<p>If you haven't installed the app yet, <a href="https://play.google.com/store/apps/details?id=com.bluenv.goassessment">click here to download</a>.</p>
			</body>
		</html>
	`)
	})
	app.Get("/routes", func(c *fiber.Ctx) error {
		routes := app.Stack()
		uniqueRoutes := make(map[string]bool)
		var result string

		for _, handlers := range routes {
			for _, route := range handlers {
				routeKey := fmt.Sprintf(" %s", route.Path)
				if !uniqueRoutes[routeKey] {
					uniqueRoutes[routeKey] = true
					result += routeKey + "\n"
				}
			}
		}

		return c.SendString(result)
	})

	return app
}
