package svc

import (
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/repo"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type TaskService interface {
	Create(task entity.CreateTaskRequestDto) error
	List(filter map[string]string, page, limit int) (*response.PaginateResponse, error)
	Index(taskId int32) (entity.TaskDetailResponseDto, error)
	Delete(taskId int32) error
	Update(taskId int32, task entity.UpdateTaskRequestDto) error
}

type TaskServiceImpl struct {
	repo repo.TaskRepository
}

func NewTaskService(repo repo.TaskRepository) TaskService {
	return &TaskServiceImpl{repo: repo}
}

func (t *TaskServiceImpl) Update(taskId int32, task entity.UpdateTaskRequestDto) error {
	return t.repo.Update(entity.Task{
		ID:          taskId,
		Title:       task.Title,
		Description: task.Description,
		Content:     task.Content,
	})

}

func (t *TaskServiceImpl) List(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	return t.repo.List(filter, page, limit)
}

func (t *TaskServiceImpl) Index(taskId int32) (entity.TaskDetailResponseDto, error) {
	return t.repo.Index(taskId)
}

func (t *TaskServiceImpl) Create(task entity.CreateTaskRequestDto) error {
	taskEntity := entity.Task{
		Title:       task.Title,
		Description: task.Description,
		Content:     task.Content,
		AttachedURL: task.AttachedURL,
	}
	return t.repo.Create(taskEntity)
}

func (t *TaskServiceImpl) Delete(taskId int32) error {
	return t.repo.Delete(taskId)

}
