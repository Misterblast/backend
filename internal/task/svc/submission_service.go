package svc

import (
	"fmt"
	"mime/multipart"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	"github.com/ghulammuzz/misterblast/pkg/agent"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type TaskSubmissionService interface {
	SubmitTask(taskId int64, userId int64, dto entity.SubmitTaskRequestDto) error
	GiveScore(submissionId int64, dto entity.ScoreSubmissionRequestDto) error
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
	err := s.repo.Create(taskId, userId, dto.Answer, "")
	if err != nil {
		log.Error("[TaskSubmissionSvc] Failed to create task submission", "error", err)
		return app.NewAppError(500, "failed to create task submission")
	}

	if dto.AttachedURL != nil {
		go func(file *multipart.FileHeader, taskId int64, userId int64) {
			url, err := agent.FileUploadProxyRESTY(file, fmt.Sprintf("/prod/user/%d/task-submission/%d", taskId, userId))
			if err != nil {
				log.Error("[TaskSubmissionSvc] Failed to upload attachment in background", "error", err)
				return
			}
			err = s.repo.UpdateAttachmentURL(taskId, userId, url)
			if err != nil {
				log.Error("[TaskSubmissionSvc] Failed to update attachment URL after upload", "error", err)
			} else {
				log.Info("[TaskSubmissionSvc] Successfully updated attachment URL", "url", url)
			}
		}(dto.AttachedURL, taskId, userId)
	}

	return nil
}

func (s *TaskSubmissionServiceImpl) GiveScore(submissionId int64, dto entity.ScoreSubmissionRequestDto) error {
	if dto.Score < 0 || dto.Score > 100 {
		return app.NewAppError(400, "score must be between 0 and 100 : task.submission.score_invalid")
	}
	return s.repo.ScoreSubmission(submissionId, dto)
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
