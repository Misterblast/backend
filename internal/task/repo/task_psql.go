package repo

import (
	"database/sql"
	"fmt"
	"github.com/ghulammuzz/misterblast/internal/models"
	entity2 "github.com/ghulammuzz/misterblast/internal/storage/entity"
	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"net/http"
	"strings"
)

type TaskRepository interface {
	List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error)
	ListAttachments(taskId int32) ([]entity2.Attachment, error)
	Create(task entity.Task) (int64, error)
	Index(taskId int32) (entity.TaskDetailResponseDto, error)
	Update(task entity.Task) error
	Delete(taskId int32) error
	InsertAttachments(taskId int64, attachmentIds []int64) error
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
		log.Error("[Repo][Tasks] failed to query count, cause : %s", err.Error())
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
			log.Error("[Repo][Tasks] failed to scan tasks, cause : %s", err.Error())
			return response, app.NewAppError(http.StatusInternalServerError, "failed to scan task")
		}
		tasks = append(tasks, task)
	}

	response.Items = tasks
	return response, nil
}

func (r *TaskRepositoryImpl) Index(taskId int32) (entity.TaskDetailResponseDto, error) {
	var task entity.TaskDetailResponseDto
	tasksQuery := `SELECT 
    t.id, t.title, t.description, t.content, t.updated_at, 
    u.id, u.name, u.email, COALESCE(u.img_url, ''), 
    u2.id, u2.name, u2.email, COALESCE(u2.img_url, '')  
	FROM tasks t 
	    LEFT JOIN users u ON t.author = u.id 
	    LEFT JOIN users u2 ON t.updated_by = u2.id 
	WHERE t.id = $1 AND t.deleted_at IS NULL`

	row := r.db.QueryRow(tasksQuery, taskId)
	if err := row.Scan(
		&task.ID, &task.Title, &task.Description, &task.Content, &task.LastUpdatedAt,
		&task.Author.ID, &task.Author.Name, &task.Author.Email, &task.Author.ImgUrl,
		&task.LastUpdatedBy.ID, &task.LastUpdatedBy.Name, &task.LastUpdatedBy.Email, &task.LastUpdatedBy.ImgUrl,
	); err != nil {
		log.Error("[Repo][Tasks] failed to scan tasks, cause : %s", err.Error())
		if err.Error() == "sql: no rows in result set" {
			return task, app.NewAppError(http.StatusNotFound, "task not found")
		}
		return task, app.NewAppError(http.StatusInternalServerError, "failed to scan task")
	}

	attachments, err := r.ListAttachments(taskId)
	if err != nil {
		return task, err
	}
	task.Attachments = attachments

	return task, nil
}

func (r *TaskRepositoryImpl) Create(task entity.Task) (int64, error) {
	query := "INSERT INTO tasks (title, description, content, author, updated_by) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	var taskId int64
	err := r.db.QueryRow(query, task.Title, task.Description, task.Content, task.Author, task.Author).Scan(&taskId)
	if err != nil {
		log.Error("[Repo.Task.Create] failed to insert task, cause : %s", err.Error())
		return 0, err
	}

	return taskId, nil
}

func (r *TaskRepositoryImpl) Update(task entity.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, content = $3, updated_by = $4, updated_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $5`
	res, err := r.db.Exec(query, task.Title, task.Description, task.Content, task.UpdatedBy, task.ID)
	if err != nil {
		log.Error("[Repo.Task.Update] failed to update task, cause: %s", err.Error())
		return app.NewAppError(http.StatusInternalServerError, "failed to update task")
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return app.NewAppError(http.StatusNotFound, "task not found")
	}

	return nil
}

func (r *TaskRepositoryImpl) Delete(taskId int32) error {
	query := `UPDATE tasks SET deleted_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $1;`
	_, err := r.db.Exec(query, taskId)
	if err != nil {
		log.Error("[Repo][Tasks] failed to delete task, cause : %s", err.Error())
		return app.NewAppError(http.StatusInternalServerError, "failed to delete task")
	}
	return nil
}

func (r *TaskRepositoryImpl) ListAttachments(taskId int32) ([]entity2.Attachment, error) {
	var attachments []entity2.Attachment
	query := "SELECT a.type, a.url FROM task_attachments LEFT JOIN attachments a ON task_attachments.attachment_id = a.id WHERE task_id = $1"
	rows, err := r.db.Query(query, taskId)
	if err != nil {
		log.Error("[Repo][Tasks] failed to query attachments, cause : %s", err.Error())
		return attachments, app.NewAppError(http.StatusInternalServerError, "failed to get attachments")
	}
	defer rows.Close()
	for rows.Next() {
		var attachment entity2.Attachment
		if err := rows.Scan(&attachment.Type, &attachment.Url); err != nil {
			log.Error("[Repo][Tasks] failed to scan attachments, cause : %s", err.Error())
			return attachments, app.NewAppError(http.StatusInternalServerError, "failed to scan attachment")
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

func (r *TaskRepositoryImpl) InsertAttachments(taskId int64, attachmentIds []int64) error {
	query := "INSERT INTO task_attachments (task_id, attachment_id) VALUES "
	var values []interface{}
	var placeholders []string

	for i, attachmentId := range attachmentIds {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		values = append(values, taskId, attachmentId)
	}

	query += strings.Join(placeholders, ", ")

	if _, err := r.db.Exec(query, values...); err != nil {
		log.Error("[Repo.Task.Create] failed to insert task_attachments, cause: %s", err.Error())
		return app.NewAppError(http.StatusInternalServerError, "failed to insert task attachments")
	}

	return nil
}
