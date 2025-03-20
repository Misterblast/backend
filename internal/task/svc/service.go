package service

import (
	"errors"
	"mime/multipart"

	"github.com/ghulammuzz/misterblast/helper"
	storageEntity "github.com/ghulammuzz/misterblast/internal/storage/entity"
	"github.com/ghulammuzz/misterblast/internal/storage/svc"
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log2 "github.com/gofiber/fiber/v2/log"
)

type TaskService interface {
	Create(userId int32, task entity.CreateTaskRequestDto) error
	List() (entity.ListTaskResponseDto, error)
}

type TaskServiceImpl struct {
	repo           repo.TaskRepository
	storageService svc.StorageService
}

func NewTaskService(repo repo.TaskRepository, storageService svc.StorageService) TaskService {
	return &TaskServiceImpl{repo: repo, storageService: storageService}
}

func (t *TaskServiceImpl) List() (entity.ListTaskResponseDto, error) {
	return entity.ListTaskResponseDto{}, nil
}

func (t *TaskServiceImpl) Create(userId int32, task entity.CreateTaskRequestDto) error {
	attachments, err := t.uploadFilesToStorageService(task.Attachments)
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			return app.NewAppError(500, "failed to upload file")
		}
		return appErr
	}
	taskEntity := entity.Task{
		Title:       task.Title,
		Description: task.Description,
		Content:     task.Content,
		Author:      userId,
		Attachments: attachments,
	}
	if err := t.repo.Create(taskEntity); err != nil {
		log2.Errorf("[Svc.Task.Create] failed to create task, cause: %v", err)
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			return app.NewAppError(500, "failed to create task")
		}
		return err
	}
	return nil
}

func (t *TaskServiceImpl) uploadFilesToStorageService(files []*multipart.FileHeader) ([]entity.TaskAttachment, error) {
	var attachments []entity.TaskAttachment

	for _, attachment := range files {
		if !helper.ValidateFileSize(attachment, 10*1024*1024) {
			return []entity.TaskAttachment{}, app.NewAppError(400, "file size is too large")
		}
		response, err := t.storageService.UploadFile(storageEntity.UploadFileRequestDto{
			Key:  attachment.Filename,
			File: attachment,
		})
		if err != nil {
			var appErr *app.AppError
			if !errors.As(err, &appErr) {
				return []entity.TaskAttachment{}, app.NewAppError(500, "failed to upload file")
			}
			return []entity.TaskAttachment{}, appErr
		}
		attachments = append(attachments, entity.TaskAttachment{
			Type: helper.GetFileExtension(attachment),
			Url:  response.Data.Url,
		})
	}
	return attachments, nil
}
