package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mlog "log/slog"

	config "github.com/ghulammuzz/misterblast/config/postgres"
	"github.com/ghulammuzz/misterblast/config/validator"
	class "github.com/ghulammuzz/misterblast/internal/class/di"
	email "github.com/ghulammuzz/misterblast/internal/email/di"
	"github.com/ghulammuzz/misterblast/internal/health"
	lesson "github.com/ghulammuzz/misterblast/internal/lesson/di"
	question "github.com/ghulammuzz/misterblast/internal/question/di"
	quiz "github.com/ghulammuzz/misterblast/internal/quiz/di"
	set "github.com/ghulammuzz/misterblast/internal/set/di"
	task "github.com/ghulammuzz/misterblast/internal/task/di"
	user "github.com/ghulammuzz/misterblast/internal/user/di"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ghulammuzz/misterblast/pkg/log"
	metrics "github.com/ghulammuzz/misterblast/pkg/prom"
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

	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()
		path := c.Path()
		method := c.Method()
		status := fmt.Sprintf("%d", c.Response().StatusCode())

		metrics.RequestCounter.WithLabelValues(path, method).Inc()
		metrics.RequestDuration.WithLabelValues(path, method).Observe(duration)

		if err != nil {
			metrics.ErrorCounter.WithLabelValues(path, method, status).Inc()
		}

		return err
	})

	log.Info(os.Getenv("JWT_SECRET"))
	app.Get("/hc", health.HealthCheck(db))
	// app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
	http.Handle("/metrics", promhttp.Handler())
	promServer := &http.Server{
		Addr: ":3002",
	}

	api := app.Group("/api")
	class.InitializedClassService(db).Router(api)
	lesson.InitializedLessonService(db, validator.Validate).Router(api)
	set.InitializedSetService(db, validator.Validate).Router(api)
	question.InitializedQuestionService(db, validator.Validate).Router(api)
	user.InitializedUserService(db, validator.Validate).Router(api)
	email.InitializedEmailService(db, validator.Validate).Router(api)
	quiz.InitializedQuizService(db, validator.Validate).Router(api)
	task.InitializeTaskService(db, validator.Validate).Router(api)

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
	go func() {
		log.Info("Listening and serving prometheus exporter", mlog.Int("port", 3002))
		if err := promServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed listening prometheus exporter", mlog.Any("err", err))
			panic(err)
		}
	}()
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
			log.Error("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Server is shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Info("Server forced to shutdown: %v", err)
	}
	log.Info("Server exited properly")

}
