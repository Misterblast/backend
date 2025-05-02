package di

import (
	"database/sql"

	quizHandler "github.com/ghulammuzz/misterblast/internal/quiz/handler"
	quizRepo "github.com/ghulammuzz/misterblast/internal/quiz/repo"
	quizSvc "github.com/ghulammuzz/misterblast/internal/quiz/svc"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedQuizServiceFake(sb *sql.DB, val *validator.Validate) *quizHandler.QuizHandler {
	wire.Build(
		quizHandler.NewQuizHandler,
		quizSvc.NewQuizService,
		quizRepo.NewQuizRepository,
	)

	return &quizHandler.QuizHandler{}
}
