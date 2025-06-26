package repo

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/ghulammuzz/misterblast/pkg/sqlutils"
)

type TaskSubmissionRepository interface {
	Create(taskId int64, userId int64, answer string, attachedURL string) error
	ScoreSubmission(submissionId int64, submissionDto entity.ScoreSubmissionRequestDto) error
	UpdateAttachmentURL(taskId int64, userId int64, url string) error

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

func (t *TaskSubmissionRepositoryImpl) Create(taskId int64, userId int64, answer string, attachedURL string) error {
	query := `
		INSERT INTO public.task_submissions (task_id, user_id, answer, attachment_url)
		VALUES ($1, $2, $3, $4)
	`
	_, err := t.db.Exec(query, taskId, userId, answer, attachedURL)
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
		log.Error("[TaskSubmissionRepo] failed to query submissions by task ID, cause: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var submissions []entity.TaskListSubmissionResponseDto
	for rows.Next() {
		var s entity.TaskListSubmissionResponseDto
		err = rows.Scan(&s.ID, &s.Title, &s.Description, &s.Content, &s.SubmittedAt, &s.ScoredAt, &s.Feedback, &s.Score)
		if err != nil {
			log.Error("[TaskSubmissionRepo] failed to scan submission, cause: %s", err.Error())
			return nil, err
		}
		submissions = append(submissions, s)
	}

	var total int64
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM task_submissions ts %s`, where)
	err = t.db.QueryRow(countQuery, taskId).Scan(&total)
	if err != nil {
		log.Error("[TaskSubmissionRepo] failed to count submissions, cause: %s", err.Error())
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
	} else if filter["type"] == "old" {
		where += " AND to_timestamp(ts.created_at) < now() - interval '7 days'"
		order = "ORDER BY ts.created_at ASC"
	} else {
		order = "ORDER BY ts.created_at DESC"
	}

	query := fmt.Sprintf(`
		SELECT  
			ts.id, 
			t.title,
			ts.answer, 
			ts.attachment_url, 
			t.description, 
			t.content, 
			ts.created_at, 
			ts.scored_at, 
			ts.feedback, 
			ts.score 
		FROM task_submissions ts
		JOIN tasks t ON t.id = ts.task_id
		%s
		%s
		LIMIT $2 OFFSET $3
	`, where, order)

	rows, err := t.db.Query(query, userId, limit, offset)
	if err != nil {
		log.Error("[TaskSubmissionRepo] failed to query submissions by user ID, cause: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	var submissions []entity.TaskListSubmissionResponseDto
	for rows.Next() {
		var s entity.TaskListSubmissionResponseDto

		var attachmentURL sql.NullString
		var scoredAt sql.NullInt64
		var feedback sql.NullString
		var score sql.NullInt32

		err = rows.Scan(
			&s.ID,
			&s.Title,
			&s.Answer,
			&attachmentURL,
			&s.Description,
			&s.Content,
			&s.SubmittedAt,
			&scoredAt,
			&feedback,
			&score,
		)
		if err != nil {
			log.Error("[TaskSubmissionRepo] failed to scan row, cause: %s", err.Error())
			return nil, err
		}

		s.AttachedURL = sqlutils.ToString(attachmentURL)
		s.Feedback = sqlutils.ToString(feedback)
		s.Score = sqlutils.ToInt32(score)
		s.ScoredAt = sqlutils.ToInt64(scoredAt)

		submissions = append(submissions, s)
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM task_submissions ts %s`, where)
	var total int64
	err = t.db.QueryRow(countQuery, userId).Scan(&total)
	if err != nil {
		log.Error("[TaskSubmissionRepo] failed to count submissions, cause: %s", err.Error())
		return nil, err
	}

	return &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  submissions,
	}, nil
}

func (t *TaskSubmissionRepositoryImpl) ScoreSubmission(submissionId int64, submissionDto entity.ScoreSubmissionRequestDto) error {
	query := `
		UPDATE public.task_submissions
		SET score = $1, feedback = $2, scored_at = EXTRACT(EPOCH FROM now())
		WHERE id = $3
	`
	res, err := t.db.Exec(query, submissionDto.Score, submissionDto.Feedback, submissionId)
	if err != nil {
		log.Error("[TaskSubmissionRepo] failed to update submission score, cause: %s", err.Error())
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

	var taskAttachmentURL, answerAttachmentURL, feedback sql.NullString
	var score sql.NullInt64
	var scoredAt sql.NullInt64
	var submittedAt sql.NullInt64

	query := `
		SELECT
			ts.id,
			COALESCE(t.title, ''),
			COALESCE(t.description, ''),
			COALESCE(t.content, ''),
			t.attachment_url,
			ts.answer,
			ts.attachment_url,
			ts.score,
			ts.scored_at,
			ts.created_at,
			ts.feedback
		FROM task_submissions ts 
		LEFT JOIN tasks t ON t.id = ts.task_id
		WHERE ts.id = $1
	`

	err := t.db.QueryRow(query, submissionId).Scan(
		&response.ID,
		&response.Title,
		&response.Description,
		&response.Content,
		&taskAttachmentURL,
		&response.Answer,
		&answerAttachmentURL,
		&score,
		&scoredAt,
		&submittedAt,
		&feedback,
	)
	if err != nil {
		log.Error("[TaskSubmissionRepo] failed to scan submission detail, cause : %s", err.Error())
		return nil, err
	}

	response.TaskAttachmentUrl = sqlutils.ToString(taskAttachmentURL)
	response.AnswerAttachmentUrl = sqlutils.ToString(answerAttachmentURL)
	response.Score = sqlutils.ToInt64(score)
	response.ScoredAt = sqlutils.ToInt64(scoredAt)
	response.SubmittedAt = sqlutils.ToInt64(submittedAt)
	response.Feedback = sqlutils.ToString(feedback)

	return &response, nil
}

func (t *TaskSubmissionRepositoryImpl) UpdateAttachmentURL(taskId int64, userId int64, url string) error {
	query := `
		UPDATE task_submissions
		SET attachment_url = $1
		WHERE task_id = $2 AND user_id = $3
		`
	_, err := t.db.Exec(query, url, taskId, userId)
	if err != nil {
		log.Error("[TaskSubmissionRepo] Failed to update attachment URL", "error", err, "taskId", taskId, "userId", userId)
		return err
	}

	log.Info("[TaskSubmissionRepo] Successfully updated attachment URL", "taskId", taskId, "userId", userId, "url", url)
	return nil
}

func NewTaskSubmissionRepository(db *sql.DB) TaskSubmissionRepository {
	return &TaskSubmissionRepositoryImpl{db: db}
}
