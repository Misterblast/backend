//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"
	repo2 "github.com/ghulammuzz/misterblast/internal/storage/repo"

	"github.com/ghulammuzz/misterblast/internal/storage/svc"
	"github.com/ghulammuzz/misterblast/internal/task/handler"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	service "github.com/ghulammuzz/misterblast/internal/task/svc"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializeTaskService(sb *sql.DB, val *validator.Validate, storageService svc.StorageService, storageRepo repo2.StorageRepository) *handler.TaskHandler {
	wire.Build(
		handler.NewTaskHandler,
		service.NewTaskService,
		repo.NewTaskRepository,
	)

	return &handler.TaskHandler{}
}
