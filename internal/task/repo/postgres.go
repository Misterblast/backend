package repo

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/ghulammuzz/misterblast/internal/models"
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type TaskRepository interface {
	List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error)
	Create(task entity.Task) error
	Update(task entity.Task) error
	Delete(taskId int32) error
}

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{db: db}
}
func (r *TaskRepositoryImpl) List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error) {
	var response models.PaginationResponse[entity.TaskResponseDto]

	queryCount := "SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL"
	var args []interface{}

	if request.Search != "" {
		queryCount += " AND title ILIKE $1"
		args = append(args, "%"+request.Search+"%")
	}

	err := r.db.QueryRow(queryCount, args...).Scan(&response.Total)
	if err != nil {
		log.Error("[Repo][Tasks] failed to query count, cause: %v", err) // Corrected log.Error -> log.Errorf
		return response, app.NewAppError(http.StatusInternalServerError, "failed to get count")
	}
	query := `SELECT 
				t.id, t.title, t.description, t.content, t.updated_at,
				u.id, u.name, u.email, COALESCE(u.img_url, ''),
				u2.id, u2.name, u2.email, COALESCE(u2.img_url, '')
			  FROM tasks t 
			  LEFT JOIN users u ON t.author = u.id
			  LEFT JOIN users u2 ON updated_by = u2.id
			  WHERE t.deleted_at IS NULL`

	args = []interface{}{}
	argIndex := 1

	if request.Search != "" {
		query += fmt.Sprintf(" AND t.title ILIKE $%d", argIndex)
		args = append(args, "%"+request.Search+"%")
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY t.updated_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, 10, (request.Page-1)*10)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Error("[Repo][Tasks] failed to query tasks, cause : %v", err)
		return response, app.NewAppError(http.StatusInternalServerError, "failed to get tasks")
	}
	defer rows.Close()
	tasks := make([]entity.TaskResponseDto, 0)
	for rows.Next() {
		var task entity.TaskResponseDto
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Content, &task.LastUpdatedAt,
			&task.Author.ID, &task.Author.Name, &task.Author.Email, &task.Author.ImgUrl,
			&task.LastUpdatedBy.ID, &task.LastUpdatedBy.Name, &task.LastUpdatedBy.Email, &task.LastUpdatedBy.ImgUrl,
		); err != nil {
			log.Error("[Repo][Tasks] failed to scan tasks, cause : %v", err)
			return response, app.NewAppError(http.StatusInternalServerError, "failed to scan task")
		}
		tasks = append(tasks, task)
	}

	response.Items = tasks
	return response, nil
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
		log.Error("[Repo.Task.Create] failed to insert task, cause: %v", err)
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
			log.Error("[Repo.Task.Create] failed to insert attachments, cause: %v", err)
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

//if len(tasks) > 0 {
//		var taskIds []int32
//		for _, task := range tasks {
//			taskIds = append(taskIds, task.ID)
//		}
//
//		queryAttachments := "SELECT task_id, type, url FROM task_attachments WHERE task_id IN ($1)"
//		rows, err := r.db.Query(queryAttachments, taskIds)
//		if err != nil {
//			return models.PaginationResponse[entity.TaskResponseDto]{}, app.NewAppError(http.StatusInternalServerError, "failed to get attachments")
//		}
//		defer rows.Close()
//		var attachments []entity.TaskAttachment
//		for rows.Next() {
//			var attachment entity.TaskAttachment
//			if err := rows.Scan(&attachment.TaskId, &attachment.Type, &attachment.Url); err != nil {
//				return models.PaginationResponse[entity.TaskResponseDto]{}, app.NewAppError(http.StatusInternalServerError, "failed to scan attachment")
//			}
//			attachments = append(attachments, attachment)
//		}
//	}
