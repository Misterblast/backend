package di

import (
	"database/sql"

	userHandler "github.com/ghulammuzz/misterblast/internal/user/handler"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"
	userSvc "github.com/ghulammuzz/misterblast/internal/user/svc"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedUserServiceFake(sb *sql.DB, val *validator.Validate) *userHandler.UserHandler {
	wire.Build(
		userHandler.NewUserHandler,
		userSvc.NewUserService,
		userRepo.NewUserRepository,
	)

	return &userHandler.UserHandler{}
}
