package service

import (
	"errors"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/ghulammuzz/misterblast/helper"
	"github.com/ghulammuzz/misterblast/internal/models"
	storageEntity "github.com/ghulammuzz/misterblast/internal/storage/entity"
	"github.com/ghulammuzz/misterblast/internal/storage/svc"
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/gofiber/fiber/v2"
	log2 "github.com/gofiber/fiber/v2/log"
)

type TaskService interface {
	Create(userId int32, task entity.CreateTaskRequestDto) error
	List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error)
	Index(taskId int32) (entity.TaskDetailResponseDto, error)
}

type TaskServiceImpl struct {
	repo           repo.TaskRepository
	storageService svc.StorageService
}

func NewTaskService(repo repo.TaskRepository, storageService svc.StorageService) TaskService {
	return &TaskServiceImpl{repo: repo, storageService: storageService}
}

func (t *TaskServiceImpl) List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error) {
	var response models.PaginationResponse[entity.TaskResponseDto]
	response.Limit = 10
	response.Page = int(request.Page)
	queryResponse, err := t.repo.List(request)
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			return models.PaginationResponse[entity.TaskResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "failed to list tasks")
		}
		return models.PaginationResponse[entity.TaskResponseDto]{}, appErr
	}
	response.Total = queryResponse.Total
	for i := range queryResponse.Items {
		timeUnix, _ := strconv.ParseInt(queryResponse.Items[i].LastUpdatedAt, 10, 64)
		queryResponse.Items[i].LastUpdatedAt = time.Unix(timeUnix, 0).Format("2006-01-02 15:04:05")
	}
	response.Items = queryResponse.Items

	return response, nil
}

func (t *TaskServiceImpl) Index(taskId int32) (entity.TaskDetailResponseDto, error) {
	panic("unimplemented")
}

func (t *TaskServiceImpl) Create(userId int32, task entity.CreateTaskRequestDto) error {
	attachments, err := t.uploadFilesToStorageService(task.Attachments)
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			return app.NewAppError(fiber.StatusInternalServerError, "failed to upload file")
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
			return app.NewAppError(fiber.StatusInternalServerError, "failed to create task")
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
				return []entity.TaskAttachment{}, app.NewAppError(fiber.StatusInternalServerError, "failed to upload file")
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
