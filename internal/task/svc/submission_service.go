package svc

import (
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type TaskSubmissionService interface {
	SubmitTask(taskId int64, userId int64, dto entity.SubmitTaskRequestDto) error
	GiveScore(submissionId int64, userId int64, dto entity.ScoreSubmissionRequestDto) error
	GetSubmissionsByUser(filter map[string]string, userId int64) (*response.PaginateResponse, error)
	GetSubmissionsByTask(filter map[string]string, taskId int64) (*response.PaginateResponse, error)
	GetSubmissionDetailById(submissionId int64) (*entity.TaskSubmissionDetailResponseDto, error)
}

type TaskSubmissionServiceImpl struct {
	repo repo.TaskSubmissionRepository
}

func NewTaskSubmissionService(repo repo.TaskSubmissionRepository) TaskSubmissionService {
	return &TaskSubmissionServiceImpl{repo: repo}
}

func (s *TaskSubmissionServiceImpl) SubmitTask(taskId int64, userId int64, dto entity.SubmitTaskRequestDto) error {
	if dto.Answer == "" {
		return app.NewAppError(400, "answer is required : task.submission.answer_required")
	}
	return s.repo.Create(taskId, userId, dto)
}

func (s *TaskSubmissionServiceImpl) GiveScore(submissionId int64, userId int64, dto entity.ScoreSubmissionRequestDto) error {
	if dto.Score < 0 || dto.Score > 100 {
		return app.NewAppError(400, "score must be between 0 and 100 : task.submission.score_invalid")
	}
	return s.repo.ScoreSubmission(submissionId, userId, dto)
}

func (s *TaskSubmissionServiceImpl) GetSubmissionsByUser(filter map[string]string, userId int64) (*response.PaginateResponse, error) {
	return s.repo.ListByUserId(filter, userId)
}

func (s *TaskSubmissionServiceImpl) GetSubmissionsByTask(filter map[string]string, taskId int64) (*response.PaginateResponse, error) {
	return s.repo.LIstByTaskId(filter, taskId)
}

func (s *TaskSubmissionServiceImpl) GetSubmissionDetailById(submissionId int64) (*entity.TaskSubmissionDetailResponseDto, error) {
	return s.repo.SubmissionDetailById(submissionId)
}
