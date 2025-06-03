package app

import (
	"database/sql"
	"fmt"

	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/gofiber/fiber/v2"
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
