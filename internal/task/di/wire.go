//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"

	"github.com/ghulammuzz/misterblast/internal/task/handler"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	service "github.com/ghulammuzz/misterblast/internal/task/svc"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializeTaskServiceFake(sb *sql.DB, val *validator.Validate) *handler.TaskHandler {
	wire.Build(
		handler.NewTaskHandler,
		service.NewTaskService,
		repo.NewTaskRepository,
	)

	return &handler.TaskHandler{}
}

func InitializeTaskSubmissionServiceFake(sb *sql.DB, val *validator.Validate) *handler.TaskSubmissionHandler {
	wire.Build(
		handler.NewTaskSubmissionHandler,
		service.NewTaskSubmissionService,
		repo.NewTaskSubmissionRepository,
	)

	return &handler.TaskSubmissionHandler{}
}
