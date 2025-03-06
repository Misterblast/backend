package di

import (
	"database/sql"

	setHandler "github.com/ghulammuzz/misterblast/internal/set/handler"
	setRepo "github.com/ghulammuzz/misterblast/internal/set/repo"
	setSvc "github.com/ghulammuzz/misterblast/internal/set/svc"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedSetServiceFake(sb *sql.DB, val *validator.Validate) *setHandler.SetHandler {
	wire.Build(
		setHandler.NewSetHandler,
		setSvc.NewSetService,
		setRepo.NewSetRepository,
	)

	return &setHandler.SetHandler{}
}
