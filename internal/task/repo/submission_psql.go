package repo

import (
	"database/sql"
	"net/http"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type TaskSubmissionRepository interface {
	Create(taskId int64, userId int64, answer string) error
	ScoreSubmission(submissionId int64, userId int64, submissionDto entity.ScoreSubmissionRequestDto) error
}

type TaskSubmissionRepositoryImpl struct {
	db *sql.DB
}

func (r *TaskSubmissionRepositoryImpl) Create(taskId int64, userId int64, answer string) error {
	query := "INSERT INTO task_submissions (task_id, user_id, answer) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, taskId, userId, answer)
	if err != nil {
		log.Error("[Repo.TaskSubmissions.InsertAttachments] failed to insert task submission, cause: %s", err.Error())
		return app.NewAppError(http.StatusInternalServerError, "failed to submit task")
	}
	return nil
}

func (r *TaskSubmissionRepositoryImpl) ScoreSubmission(submissionId int64, userId int64, submissionDto entity.ScoreSubmissionRequestDto) error {
	query := "UPDATE task_submissions SET score = $1, feedback = $2, scored_by = $3, scored_at = NOW() WHERE id = $4"
	res, err := r.db.Exec(query, submissionDto.Score, submissionDto.Feedback, userId, submissionId)
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return app.NewAppError(http.StatusNotFound, "task submission not found")
	}
	if err != nil {
		log.Error("[Repo.TaskSubmissions.ScoreSubmission] failed to update task submission score, cause: %s", err.Error())
		return app.NewAppError(http.StatusInternalServerError, "failed to update task submission score")
	}
	return nil
}