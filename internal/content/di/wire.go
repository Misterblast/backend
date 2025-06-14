package di

import (
	"database/sql"

	contentHandler "github.com/ghulammuzz/misterblast/internal/content/handler"
	contentRepo "github.com/ghulammuzz/misterblast/internal/content/repo"
	contentSvc "github.com/ghulammuzz/misterblast/internal/content/svc"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func InitializedContentServiceFake(db *sql.DB, redis *redis.Client, val *validator.Validate) *contentHandler.ContentHandler {
	wire.Build(
		contentHandler.NewContentHandler,
		contentSvc.NewContentService,
		contentRepo.NewContentRepository,
	)

	return &contentHandler.ContentHandler{}
}

func InitializedAuthorServiceFake(db *sql.DB, redis *redis.Client, val *validator.Validate) *contentHandler.AuthorHandler {
	wire.Build(
		contentHandler.NewAuthorHandler,
		contentSvc.NewAuthorService,
		contentRepo.NewAuthorRepository,
	)

	return &contentHandler.AuthorHandler{}
}
