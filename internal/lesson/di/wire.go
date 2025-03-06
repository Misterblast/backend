package di

import (
	"database/sql"

	lessonHandler "github.com/ghulammuzz/misterblast/internal/lesson/handler"
	lessonRepo "github.com/ghulammuzz/misterblast/internal/lesson/repo"
	lessonSvc "github.com/ghulammuzz/misterblast/internal/lesson/svc"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedLessonServiceFake(sb *sql.DB, val *validator.Validate) *lessonHandler.LessonHandler {
	wire.Build(
		lessonHandler.NewLessonHandler,
		lessonSvc.NewLessonService,
		lessonRepo.NewLessonRepository,
	)

	return &lessonHandler.LessonHandler{}
}
