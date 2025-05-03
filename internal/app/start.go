package app

import (
	"fmt"
	"os"

	pg "github.com/ghulammuzz/misterblast/config/postgres"
	cache "github.com/ghulammuzz/misterblast/config/redis"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

func Start() {
	db, _ := pg.InitPostgres()
	defer db.Close()

	redis, _ := cache.InitRedis()
	defer redis.Close()

	app := SetupRouter(db, redis)

	RegisterHealthRoutes(app, db)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
			log.Error("Error starting server: %v", err)
		}
	}()

	StartPrometheusExporter()

	GracefulShutdown(app)
}
