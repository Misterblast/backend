package di

import (
	"database/sql"

	classHandler "github.com/ghulammuzz/misterblast/internal/class/handler"
	classRepo "github.com/ghulammuzz/misterblast/internal/class/repo"
	classSvc "github.com/ghulammuzz/misterblast/internal/class/svc"
	"github.com/google/wire"
)

func InitializedClassServiceFake(sb *sql.DB) *classHandler.ClassHandler {
	wire.Build(
		classHandler.NewClassHandler,
		classSvc.NewClassService,
		classRepo.NewClassRepository,
	)

	return &classHandler.ClassHandler{}
}
