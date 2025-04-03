package service

import (
	"errors"
	"fmt"
	"github.com/ghulammuzz/misterblast/helper"
	"github.com/ghulammuzz/misterblast/internal/models"
	storageEntity "github.com/ghulammuzz/misterblast/internal/storage/entity"
	"github.com/ghulammuzz/misterblast/internal/storage/svc"
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/gofiber/fiber/v2"
	log2 "github.com/gofiber/fiber/v2/log"
	"mime/multipart"
)

type TaskService interface {
	Create(userId int32, task entity.CreateTaskRequestDto) error
	List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error)
	Index(taskId int32) (entity.TaskDetailResponseDto, error)
	Delete(taskId int32) error
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
			log2.Errorf("[Svc.Task.List] failed to list tasks, cause: %v", err)
			return models.PaginationResponse[entity.TaskResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "failed to list tasks")
		}
		return models.PaginationResponse[entity.TaskResponseDto]{}, appErr
	}
	response.Total = queryResponse.Total
	for i := range queryResponse.Items {
		queryResponse.Items[i].LastUpdatedAt = helper.FormatUnixTime(queryResponse.Items[i].LastUpdatedAt)
	}
	response.Items = queryResponse.Items

	return response, nil
}

func (t *TaskServiceImpl) Index(taskId int32) (entity.TaskDetailResponseDto, error) {
	task, err := t.repo.Index(taskId)
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			log2.Errorf("[Svc.Task.Index] failed to get task, cause: %v", err)
			return entity.TaskDetailResponseDto{}, app.NewAppError(fiber.StatusInternalServerError, "failed to get task")
		}
		return entity.TaskDetailResponseDto{}, appErr
	}
	task.LastUpdatedAt = helper.FormatUnixTime(task.LastUpdatedAt)
	return task, nil
}

func (t *TaskServiceImpl) Create(userId int32, task entity.CreateTaskRequestDto) error {
	taskEntity := entity.Task{
		Title:       task.Title,
		Description: task.Description,
		Content:     task.Content,
		Author:      userId,
	}
	taskId, err := t.repo.Create(taskEntity)
	if err != nil {
		log2.Errorf("[Svc.Task.Create] failed to create task, cause: %v", err)
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			return app.NewAppError(fiber.StatusInternalServerError, "failed to create task")
		}
		return err
	}
	attachments, err := t.uploadFilesToStorageService(taskId, task.Attachments)
	if err != nil {
		log2.Errorf("[Svc.Task.Create] failed to upload file, cause: %v", err)
	}

	err = t.repo.InsertAttachments(taskId, attachments)
	if err != nil {
		log2.Errorf("[Svc.Task.Create] failed to insert attachments, cause: %v", err)
	}
	return nil
}

func (t *TaskServiceImpl) uploadFilesToStorageService(taskId int64, files []*multipart.FileHeader) ([]storageEntity.Attachment, error) {

	for _, attachment := range files {
		_, err := helper.GetFileType(attachment)
		if err != nil {
			return []storageEntity.Attachment{}, app.NewAppError(fiber.StatusBadRequest, "Invalid file type")
		}
		if !helper.ValidateFileSize(attachment, 10*1024*1024) {
			return []storageEntity.Attachment{}, app.NewAppError(400, "file size is too large")
		}
	}

	var attachments []storageEntity.Attachment
	for _, attachment := range files {
		fileType, _ := helper.GetFileType(attachment)
		key := fmt.Sprintf("task/%s-task/%d", fileType, taskId)
		log2.Infof("KEY : %s", key)
		response, err := t.storageService.UploadFile(storageEntity.UploadFileRequestDto{
			Key:  key,
			File: attachment,
		})
		log2.Infof("RESPONSE : %s", response.Data.Url)
		if err != nil {
			var appErr *app.AppError
			if !errors.As(err, &appErr) {
				log2.Errorf("[Svc.Task.Create] failed to upload file, cause: %v", err)
				return []storageEntity.Attachment{}, app.NewAppError(fiber.StatusInternalServerError, "failed to upload file")
			}
			return []storageEntity.Attachment{}, appErr
		}
		attachments = append(attachments, storageEntity.Attachment{
			Type: fileType,
			Url:  response.Data.Url,
		})
	}
	return attachments, nil
}

func (t *TaskServiceImpl) Delete(taskId int32) error {
	err := t.repo.Delete(taskId)
	if err != nil {
		return app.NewAppError(fiber.StatusInternalServerError, "failed to delete task")
	}
	attachments, _ := t.repo.ListAttachments(taskId)
	for _, attachment := range attachments {
		err := t.storageService.DeleteFile(helper.GetFileKey(attachment.Url))
		if err != nil {
			log2.Errorf("[Svc.Task.Delete] failed to delete file, cause: %v", err)
		}
	}
	return nil
}
