//go:build wireinject
// +build wireinject

package di

import (
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
