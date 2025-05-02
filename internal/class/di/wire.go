package di

import (
	"database/sql"

	classHandler "github.com/ghulammuzz/misterblast/internal/class/handler"
	classRepo "github.com/ghulammuzz/misterblast/internal/class/repo"
	classSvc "github.com/ghulammuzz/misterblast/internal/class/svc"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func InitializedClassServiceFake(sb *sql.DB, redis *redis.Client) *classHandler.ClassHandler {
	wire.Build(
		classHandler.NewClassHandler,
		classSvc.NewClassService,
		classRepo.NewClassRepository,
	)

	return &classHandler.ClassHandler{}
}
