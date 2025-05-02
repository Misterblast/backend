package svc

import (
	"context"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type QuestionService interface {
	// Questions
	AddQuestion(question questionEntity.SetQuestion, lang string) error
	ListQuestions(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionExample, error)
	ListQuizQuestions(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionQuiz, error)
	DeleteQuestion(id int32) error
	DetailQuestion(ctx context.Context, id int32) (questionEntity.DetailQuestionExample, error)
	EditQuestion(id int32, question questionEntity.EditQuestion) error

	// Answer
	AddQuizAnswer(answer questionEntity.SetAnswer) error
	DeleteAnswer(id int32) error
	EditQuizAnswer(id int32, answer questionEntity.EditAnswer) error

	// Admin
	ListAdmin(ctx context.Context, filter map[string]string, page, limit int) (*response.PaginateResponse, error)
}

type questionService struct {
	repo repo.QuestionRepository
}

func NewQuestionService(repo repo.QuestionRepository) QuestionService {
	return &questionService{repo: repo}
}
func (s *questionService) AddQuizAnswer(answer questionEntity.SetAnswer) error {
	return s.repo.AddQuizAnswer(answer)
}

func (s *questionService) EditQuizAnswer(id int32, question questionEntity.EditAnswer) error {
	return s.repo.EditAnswer(id, question)
}

func (s *questionService) AddQuestion(q questionEntity.SetQuestion, lang string) error {
	exists, err := s.repo.Exists(q.SetID, q.Number)
	if err != nil {
		return err
	}
	if exists {
		return app.NewAppError(409, "question number already exists in this set")
	}
	return s.repo.Add(q, lang)
}

func (s *questionService) ListQuestions(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	return s.repo.List(ctx, filter)
}

func (s *questionService) DeleteQuestion(id int32) error {
	return s.repo.Delete(id)
}

// Quiz

func (s *questionService) ListQuizQuestions(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionQuiz, error) {
	return s.repo.ListQuizQuestions(ctx, filter)
}

// admin

func (s *questionService) ListAdmin(ctx context.Context, filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	questions, err := s.repo.ListAdmin(ctx, filter, page, limit)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *questionService) EditQuestion(id int32, question questionEntity.EditQuestion) error {
	return s.repo.Edit(id, question)
}

func (s *questionService) DeleteAnswer(id int32) error {
	return s.repo.DeleteAnswer(id)
}

func (s *questionService) DetailQuestion(ctx context.Context, id int32) (questionEntity.DetailQuestionExample, error) {
	return s.repo.Detail(ctx, id)
}
