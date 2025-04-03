package svc

import (
	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	quizRepo "github.com/ghulammuzz/misterblast/internal/quiz/repo"
)

type QuizService interface {
	SubmitQuiz(req quizEntity.QuizSubmit, setID int, userID int) error
	ListAdmin(filter map[string]string, page int, limit int) ([]quizEntity.ListQuizSubmissionAdmin, error)
	List(filter map[string]string, page int, limit int, userID int) ([]quizEntity.ListQuizSubmission, error)
	GetResult(userID int) (quizEntity.QuizExp, error)
}

type quizService struct {
	repo quizRepo.QuizRepository
}

func (s *quizService) GetResult(userID int) (quizEntity.QuizExp, error) {
	return s.repo.GetLast(userID)
}

func NewQuizService(repo quizRepo.QuizRepository) QuizService {
	return &quizService{repo: repo}
}

func (s *quizService) SubmitQuiz(req quizEntity.QuizSubmit, setID int, userID int) error {
	return s.repo.Submit(req, setID, userID)
}

func (s *quizService) ListAdmin(filter map[string]string, page int, limit int) ([]quizEntity.ListQuizSubmissionAdmin, error) {
	return s.repo.ListAdmin(filter, page, limit)
}

func (s *quizService) List(filter map[string]string, page int, limit int, userID int) ([]quizEntity.ListQuizSubmission, error) {
	return s.repo.List(filter, page, limit, userID)
}
