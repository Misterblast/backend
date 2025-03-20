package repo

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	log2 "github.com/gofiber/fiber/v2/log"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
)

type TaskRepository interface {
	Create(task entity.Task) error
	Update(task entity.Task) error
	Delete(taskId int32) error
	List(request entity.ListTaskRequestDto) (entity.ListTaskResponseDto, error)
}

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) Create(task entity.Task) error {
	tx, err := r.db.Begin()
	if err != nil {
		return app.NewAppError(http.StatusInternalServerError, "failed to start transaction")
	}

	defer tx.Rollback()
	query := "INSERT INTO tasks (title, description, content, author, updated_by) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	var taskId int64
	err = tx.QueryRow(query, task.Title, task.Description, task.Content, task.Author, task.Author).Scan(&taskId)
	if err != nil {
		log2.Errorf("[Repo.Task.Create] failed to insert task, cause: %v", err)
		tx.Rollback()
		return err
	}

	if len(task.Attachments) > 0 {
		placeholderStrings := make([]string, 0, len(task.Attachments))
		values := make([]interface{}, 0, len(task.Attachments)*3)

		var sb strings.Builder
		sb.WriteString("INSERT INTO task_attachments (task_id, type, url) VALUES ")

		placeholderIndex := 1
		for _, attachment := range task.Attachments {
			placeholderStrings = append(placeholderStrings, fmt.Sprintf("($%d,$%d,$%d)", placeholderIndex, placeholderIndex+1, placeholderIndex+2))
			values = append(values, taskId, attachment.Type, attachment.Url)
			placeholderIndex += 3
		}
		sb.WriteString(strings.Join(placeholderStrings, ","))

		if _, err = tx.Exec(sb.String(), values...); err != nil {
			log2.Errorf("[Repo.Task.Create] failed to insert attachments, cause: %v", err)
			return app.NewAppError(http.StatusInternalServerError, "failed to insert attachments")
		}
	}
	tx.Commit()

	return nil
}

func (r *TaskRepositoryImpl) Update(task entity.Task) error {
	return nil
}

// Delete implements TaskRepository.
func (r *TaskRepositoryImpl) Delete(taskId int32) error {
	panic("unimplemented")
}

// List implements TaskRepository.
func (r *TaskRepositoryImpl) List(request entity.ListTaskRequestDto) (entity.ListTaskResponseDto, error) {
	panic("unimplemented")
}
