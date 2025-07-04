package entity

import "mime/multipart"

type ListTaskRequestDto struct {
	Search string
	Page   int
	Limit  int
}

type CreateTaskRequestDto struct {
	Title       string `json:"title" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"required,min=1"`
	Content     string `json:"content" validate:"required,min=1"`
	AttachedURL string `json:"attached_url" validate:"omitempty,url"`
}
type UpdateTaskRequestDto struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Content     string `form:"content"`
	AttachedURL string `form:"attached_url"`
}

type SubmitTaskRequestDto struct {
	Answer      string                `form:"answer"`
	AttachedURL *multipart.FileHeader `form:"attached_url" validate:"omitempty"`
}

type TaskResponseDto struct {
	ID            int32            `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Content       string           `json:"content"`
	LastUpdatedAt string           `json:"last_updated_at"`
	Statistic     TaskStatisticDto `json:"statistic"`
}

type TaskDetailResponseDto struct {
	ID            int32            `json:"id"`
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	Content       string           `json:"content"`
	AttachedURL   string           `json:"attached_url"`
	Statistic     TaskStatisticDto `json:"statistic"`
	LastUpdatedAt string           `json:"last_updated_at"`
}

type TaskStatisticDto struct {
	TotalAssignment        int32   `json:"total_assignment"`
	AverageAssignmentScore float32 `json:"average_assignment_score"`
}

type TaskListSubmissionResponseDto struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Answer      string `json:"answer"`
	AttachedURL string `json:"attached_url"`
	Score       int32  `json:"score"`
	ScoredAt    int64  `json:"scored_at"`
	SubmittedAt string `json:"submitted_at"`
	Feedback    string `json:"feedback"`
}

type TaskSubmissionDetailResponseDto struct {
	ID                  int64  `json:"id"`
	Title               string `json:"title"`
	Description         string `json:"description"`
	Content             string `json:"content"`
	TaskAttachmentUrl   string `json:"task_attachment_url"`
	Answer              string `json:"answer"`
	AnswerAttachmentUrl string `json:"answer_attachment_url"`
	Score               int64  `json:"score"`
	ScoredAt            int64  `json:"scored_at"`
	SubmittedAt         int64  `json:"submitted_at"`
	Feedback            string `json:"feedback"`
}
