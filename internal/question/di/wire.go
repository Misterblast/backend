package di

import (
	"database/sql"

	questionHandler "github.com/ghulammuzz/misterblast/internal/question/handler"
	questionRepo "github.com/ghulammuzz/misterblast/internal/question/repo"
	questionSvc "github.com/ghulammuzz/misterblast/internal/question/svc"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedQuestionServiceFake(sb *sql.DB, val *validator.Validate) *questionHandler.QuestionHandler {
	wire.Build(
		questionHandler.NewQuestionHandler,
		questionSvc.NewQuestionService,
		questionRepo.NewQuestionRepository,
	)

	return &questionHandler.QuestionHandler{}
}
