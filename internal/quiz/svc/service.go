package svc

import (
	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	quizRepo "github.com/ghulammuzz/misterblast/internal/quiz/repo"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type QuizService interface {
	SubmitQuiz(req quizEntity.QuizSubmit, setID int, userID int, lang string) (int, error)
	ListAdmin(filter map[string]string, page int, limit int) (*response.PaginateResponse, error)
	List(filter map[string]string, userID int) (*response.PaginateResponse, error)
	GetResult(userID int) (quizEntity.QuizExp, error)
	GetSubmissionResult(submissionId int) (quizEntity.QuizExp, error)
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

func (s *quizService) SubmitQuiz(req quizEntity.QuizSubmit, setID int, userID int, lang string) (int, error) {
	return s.repo.Submit(req, setID, userID, lang)
}

func (s *quizService) ListAdmin(filter map[string]string, page int, limit int) (*response.PaginateResponse, error) {
	return s.repo.ListAdmin(filter, page, limit)
}

func (s *quizService) List(filter map[string]string, userID int) (*response.PaginateResponse, error) {
	return s.repo.List(filter, userID)
}
func (s *quizService) GetSubmissionResult(submissionId int) (quizEntity.QuizExp, error) {
	return s.repo.GetSubmissionDetail(submissionId)
}
