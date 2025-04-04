//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"
	"github.com/ghulammuzz/misterblast/internal/storage/repo"
	"github.com/ghulammuzz/misterblast/internal/storage/svc"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
)

func InitializeStorageService(client *fiber.Client) svc.StorageService {
	wire.Build(
		svc.NewStorageService,
	)

	return nil
}

func InitializeStorageRepository(db *sql.DB) repo.StorageRepository {
	wire.Build(
		repo.NewStorageRepository,
	)

	return nil
}
