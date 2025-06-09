package repo

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type TaskSubmissionRepository interface {
	Create(taskId int64, userId int64, attachment entity.SubmitTaskRequestDto) error
	ScoreSubmission(submissionId int64, userId int64, submissionDto entity.ScoreSubmissionRequestDto) error

	// ListByUserId(filter map[string]string, userId int64) ([]entity.TaskListSubmissionResponseDto, error)
	ListByUserId(filter map[string]string, userId int64) (*response.PaginateResponse, error)
	// filter (page, limit, type(this_week, old))

	// LIstByTaskId(filter map[string]string, taskId int64) ([]entity.TaskListSubmissionResponseDto, error)
	LIstByTaskId(filter map[string]string, taskId int64) (*response.PaginateResponse, error)
	// filter (page, limit, type(this_week, old))
	SubmissionDetailById(submissionId int64) (*entity.TaskSubmissionDetailResponseDto, error)
}
type TaskSubmissionRepositoryImpl struct {
	db *sql.DB
}

func (t *TaskSubmissionRepositoryImpl) Create(taskId int64, userId int64, attachment entity.SubmitTaskRequestDto) error {
	query := `
		INSERT INTO public.task_submissions (task_id, user_id, answer, attachment_url)
		VALUES ($1, $2, $3, $4)
	`
	_, err := t.db.Exec(query, taskId, userId, attachment.Answer, attachment.AttachedURL)
	return err
}

func (t *TaskSubmissionRepositoryImpl) LIstByTaskId(filter map[string]string, taskId int64) (*response.PaginateResponse, error) {
	page, _ := strconv.Atoi(filter["page"])
	limit, _ := strconv.Atoi(filter["limit"])
	offset := (page - 1) * limit

	where := "WHERE ts.task_id = $1"
	order := ""

	if filter["type"] == "this_week" {
		where += " AND to_timestamp(ts.created_at) >= now() - interval '7 days'"
		order = "ORDER BY ts.created_at DESC"
	}

	if filter["type"] == "old" {
		where += " AND to_timestamp(ts.created_at) < now() - interval '7 days'"
		order = "ORDER BY ts.created_at ASC"
	}

	query := fmt.Sprintf(`
		SELECT ts.id, t.title,t.description, t.content, ts.created_at, ts.scored_at, ts.feedback, ts.score 
		FROM task_submissions ts
		%s
		%s
		LIMIT $2 OFFSET $3
	`, where, order)

	rows, err := t.db.Query(query, taskId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []entity.TaskListSubmissionResponseDto
	for rows.Next() {
		var s entity.TaskListSubmissionResponseDto
		err = rows.Scan(&s.ID, &s.Title, &s.Description, &s.Content, &s.SubmittedAt, &s.ScoredAt, &s.Feedback, &s.Score)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM task_submissions ts %s`, where)
	err = t.db.QueryRow(countQuery, taskId).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  submissions,
	}, nil
}

func (t *TaskSubmissionRepositoryImpl) ListByUserId(filter map[string]string, userId int64) (*response.PaginateResponse, error) {
	page, _ := strconv.Atoi(filter["page"])
	limit, _ := strconv.Atoi(filter["limit"])
	offset := (page - 1) * limit

	where := "WHERE ts.user_id = $1"
	order := ""

	if filter["type"] == "this_week" {
		where += " AND to_timestamp(ts.created_at) >= now() - interval '7 days'"
		order = "ORDER BY ts.created_at DESC"
	}
	if filter["type"] == "old" {
		where += " AND to_timestamp(ts.created_at) < now() - interval '7 days'"
		order = "ORDER BY ts.created_at ASC"
	}

	query := fmt.Sprintf(`
		SELECT  ts.id, t.title,t.description, t.content, ts.created_at, ts.scored_at, ts.feedback, ts.score 
		FROM task_submissions ts
		JOIN tasks t ON t.id = ts.task_id
		%s
		%s
		LIMIT $2 OFFSET $3
	`, where, order)

	rows, err := t.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var submissions []entity.TaskListSubmissionResponseDto
	for rows.Next() {
		var s entity.TaskListSubmissionResponseDto
		err = rows.Scan(&s.ID, &s.Title, &s.Description, &s.Content, &s.SubmittedAt, &s.ScoredAt, &s.Feedback, &s.Score)
		if err != nil {
			return nil, err
		}
		submissions = append(submissions, s)
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM task_submissions ts %s`, where)
	err = t.db.QueryRow(countQuery, userId).Scan(&total)
	if err != nil {
		return nil, err
	}

	return &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  submissions,
	}, nil
}

func (t *TaskSubmissionRepositoryImpl) ScoreSubmission(submissionId int64, userId int64, submissionDto entity.ScoreSubmissionRequestDto) error {
	query := `
		UPDATE public.task_submissions
		SET score = $1, feedback = $2, scored_at = EXTRACT(EPOCH FROM now())
		WHERE id = $3 AND user_id = $4
	`
	res, err := t.db.Exec(query, submissionDto.Score, submissionDto.Feedback, submissionId, userId)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
func (t *TaskSubmissionRepositoryImpl) SubmissionDetailById(submissionId int64) (*entity.TaskSubmissionDetailResponseDto, error) {
	var response entity.TaskSubmissionDetailResponseDto
	query := `
			SELECT
				ts.id, t.title,t.description, t.content, t.attachment_url,
				ts.answer , ts.attachment_url,
				ts.score, ts.scored_at, ts.created_at, ts.feedback
			FROM task_submissions ts 
				LEFT JOIN tasks t ON t.id = ts.task_id
			WHERE ts.id = $1
			`
	row := t.db.QueryRow(query, submissionId)
	if err := row.Scan(
		&response.ID, &response.Title, &response.Description, &response.Content, &response.TaskAttachmentUrl,
		&response.Answer, &response.AnswerAttachmentUrl,
		&response.Score, &response.ScoredAt, &response.SubmittedAt, &response.Feedback,
	); err != nil {
		return nil, err
	}
	return &response, nil
}

func NewTaskSubmissionRepository(db *sql.DB) TaskSubmissionRepository {
	return &TaskSubmissionRepositoryImpl{db: db}
}
