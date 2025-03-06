package main

import (
	"flag"
	"fmt"
	"os"

	mlog "log/slog"

	config "github.com/ghulammuzz/misterblast/config/postgres"
	"github.com/ghulammuzz/misterblast/config/validator"
	class "github.com/ghulammuzz/misterblast/internal/class/di"
	email "github.com/ghulammuzz/misterblast/internal/email/di"
	"github.com/ghulammuzz/misterblast/internal/health"
	lesson "github.com/ghulammuzz/misterblast/internal/lesson/di"
	question "github.com/ghulammuzz/misterblast/internal/question/di"
	set "github.com/ghulammuzz/misterblast/internal/set/di"
	user "github.com/ghulammuzz/misterblast/internal/user/di"

	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	env := flag.String("env", "prod", "Environment for (stg/prod)")
	flag.Parse()

	if *env == "stg" {
		err := godotenv.Load("./stg.env")
		if err != nil {
			mlog.Error("Error loading stg.env file ")
		}
		mlog.Info("Environment: staging (stg.env loaded)")
	} else {
		mlog.Info("Environment: production (using system environment variables)")
	}

	log.InitLogger("dev", false, "")
	// log.InitLogger("prod", true, "http://localhost:3100/loki/api/v1/push")

	validator.InitValidator()
}

func main() {
	db, err := config.InitPostgres()
	if err != nil {
		log.Error("Failed to initialize database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	app.Get("/hc", health.HealthCheck(db))
	api := app.Group("/api")

	class.InitializedClassService(db).Router(api)
	lesson.InitializedLessonService(db, validator.Validate).Router(api)
	set.InitializedSetService(db, validator.Validate).Router(api)
	question.InitializedQuestionService(db, validator.Validate).Router(api)
	user.InitializedUserService(db, validator.Validate).Router(api)
	email.InitializedEmailService(db, validator.Validate).Router(api)

	if err := app.Listen(fmt.Sprint(":", os.Getenv("APP_PORT"))); err != nil {
		log.Error("Failed to start the server: %v", err)
	}
}
