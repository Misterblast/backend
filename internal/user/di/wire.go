package di

import (
	"database/sql"

	quizRepo "github.com/ghulammuzz/misterblast/internal/quiz/repo"
	taskRepo "github.com/ghulammuzz/misterblast/internal/task/repo"
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
		quizRepo.NewQuizRepository,
		taskRepo.NewTaskRepository,
	)

	return &userHandler.UserHandler{}
}
