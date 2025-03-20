package entity

import (
	"mime/multipart"

	"github.com/ghulammuzz/misterblast/internal/user/entity"
)

type ListTaskRequestDto struct {
	Search string
	Page   int32
}

type CreateTaskRequestDto struct {
	Title       string                  `form:"title" validate:"required"`
	Description string                  `form:"description" validate:"required"`
	Content     string                  `form:"content" validate:"required"`
	Attachments []*multipart.FileHeader `form:"attachments"`
}

type UpdateTaskRequestDto struct {
	Title       string                  `form:"title"`
	Description string                  `form:"description"`
	Content     string                  `form:"content"`
	Attachments []*multipart.FileHeader `form:"attachments"`
}

type TaskResponseDto struct {
	ID            int32             `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Content       string            `json:"content"`
	Author        entity.DetailUser `json:"author"`
	Statistic     TaskStatisticDto  `json:"statistic"`
	LastUpdatedBy entity.DetailUser `json:"last_updated_by"`
	LastUpdatedAt string            `json:"last_updated_at"`
}

type TaskDetailResponseDto struct {
	ID            int32             `json:"id"`
	Title         string            `json:"title"`
	Description   string            `json:"description"`
	Content       string            `json:"content"`
	Author        entity.DetailUser `json:"author"`
	Statistic     TaskStatisticDto  `json:"statistic"`
	LastUpdatedBy entity.DetailUser `json:"last_updated_by"`
	LastUpdatedAt string            `json:"last_updated_at"`
}

type TaskStatisticDto struct {
	TotalAssignment        int32   `json:"total_assignment"`
	AverageAssignmentScore float32 `json:"average_assignment_score"`
}

type ListTaskResponseDto struct {
	Tasks []TaskResponseDto `json:"tasks"`
}
