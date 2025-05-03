package repo

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type TaskRepository interface {
	// List(request entity.ListTaskRequestDto) (models.PaginationResponse[entity.TaskResponseDto], error)
	List(filter map[string]string, page, limit int) (*response.PaginateResponse, error)
	Create(task entity.Task) error
	Index(taskId int32) (entity.TaskDetailResponseDto, error)
	Update(task entity.Task) error
	Delete(taskId int32) error
}

type TaskRepositoryImpl struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) List(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	var total int64
	countArgs := []interface{}{}
	queryArgs := []interface{}{}
	argIndex := 1

	// --- COUNT QUERY ---
	queryCount := "SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL"
	if search, ok := filter["search"]; ok && search != "" {
		queryCount += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		countArgs = append(countArgs, "%"+search+"%")
	}

	err := r.db.QueryRow(queryCount, countArgs...).Scan(&total)
	if err != nil {
		log.Error("[Repo][Tasks] failed to query count, cause : %s", err.Error())
		return nil, app.NewAppError(http.StatusInternalServerError, "failed to get count")
	}

	// --- DATA QUERY ---
	query := `
		SELECT t.id, t.title, t.description, t.content, t.updated_at
		FROM tasks t
		WHERE t.deleted_at IS NULL`
	if search, ok := filter["search"]; ok && search != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		queryArgs = append(queryArgs, "%"+search+"%")
		argIndex++
	}

	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	query += fmt.Sprintf(" ORDER BY t.updated_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	queryArgs = append(queryArgs, limit, offset)

	rows, err := r.db.Query(query, queryArgs...)
	if err != nil {
		log.Error("[Repo][Tasks] failed to query tasks, cause : %v", err)
		return nil, app.NewAppError(http.StatusInternalServerError, "failed to get tasks")
	}
	defer rows.Close()

	tasks := make([]entity.TaskResponseDto, 0)
	for rows.Next() {
		var task entity.TaskResponseDto
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Content, &task.LastUpdatedAt,
		); err != nil {
			log.Error("[Repo][Tasks] failed to scan tasks, cause : %s", err.Error())
			return nil, app.NewAppError(http.StatusInternalServerError, "failed to scan task")
		}
		tasks = append(tasks, task)
	}

	return &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  tasks,
	}, nil
}

func (r *TaskRepositoryImpl) Index(taskId int32) (entity.TaskDetailResponseDto, error) {
	var task entity.TaskDetailResponseDto
	tasksQuery := `SELECT 
    t.id, t.title, t.description, t.content, t.updated_at, t.attachment_url
	FROM tasks t  
	WHERE t.id = $1 AND t.deleted_at IS NULL`

	row := r.db.QueryRow(tasksQuery, taskId)
	if err := row.Scan(
		&task.ID, &task.Title, &task.Description, &task.Content, &task.LastUpdatedAt, &task.AttachedURL,
	); err != nil {
		log.Error("[Repo][Tasks] failed to scan tasks, cause : %s", err.Error())
		if err.Error() == "sql: no rows in result set" {
			return task, app.NewAppError(http.StatusNotFound, "task not found")
		}
		return task, app.NewAppError(http.StatusInternalServerError, "failed to scan task")
	}

	return task, nil
}

// done
func (r *TaskRepositoryImpl) Create(task entity.Task) error {
	query := "INSERT INTO tasks (title, description, content, attachment_url) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, task.Title, task.Description, task.Content, task.AttachedURL)
	if err != nil {
		log.Error("[Repo.Task.Create] failed to insert task, cause : %s", err.Error())
		return err
	}
	return nil
}

func (r *TaskRepositoryImpl) Update(task entity.Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, content = $3, attachment_url = $4 updated_at = EXTRACT(EPOCH FROM NOW()) WHERE id = $5`
	res, err := r.db.Exec(query, task.Title, task.Description, task.Content, task.AttachedURL, task.ID)
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
