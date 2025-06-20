package repo

import (
	"fmt"

	"github.com/ghulammuzz/misterblast/helper"
	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"

	"github.com/ghulammuzz/misterblast/pkg/response"
)

func (r *quizRepository) ListAdmin(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	query := `
		SELECT s.id, s.set_id, s.correct, s.grade, s.submitted_at,
			   u.name AS user_name,
			   l.name AS lesson_name, c.name AS class_name
		FROM quiz_submissions s
		JOIN users u ON s.user_id = u.id
		JOIN sets a ON s.set_id = a.id
		JOIN lessons l ON a.lesson_id = l.id
		JOIN classes c ON a.class_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCounter := 1

	if lesson, exists := filter["lesson_id"]; exists {
		query += fmt.Sprintf(" AND l.id = $%d", argCounter)
		args = append(args, lesson)
		argCounter++
	}
	if class, exists := filter["class_id"]; exists {
		query += fmt.Sprintf(" AND c.id = $%d", argCounter)
		args = append(args, class)
		argCounter++
	}
	if submissionType, exists := filter["type"]; exists {
		if submissionType == "this_week" {
			query += " AND s.submitted_at >= EXTRACT(EPOCH FROM NOW() - INTERVAL '7 days')"
		} else if submissionType == "old" {
			query += " AND s.submitted_at < EXTRACT(EPOCH FROM NOW() - INTERVAL '7 days')"
		}
	}

	query += " ORDER BY s.submitted_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCounter)
		args = append(args, limit)
		argCounter++
	}
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query += fmt.Sprintf(" OFFSET $%d", argCounter)
		args = append(args, offset)
	}

	// === Count Query ===
	countQuery := `
		SELECT COUNT(*) 
		FROM quiz_submissions s
		JOIN users u ON s.user_id = u.id
		JOIN sets a ON s.set_id = a.id
		JOIN lessons l ON a.lesson_id = l.id
		JOIN classes c ON a.class_id = c.id
		WHERE 1=1
	`
	countArgs := []interface{}{}
	countArgCounter := 1

	if lesson, exists := filter["lesson_id"]; exists {
		countQuery += fmt.Sprintf(" AND l.id = $%d", countArgCounter)
		countArgs = append(countArgs, lesson)
		countArgCounter++
	}
	if class, exists := filter["class_id"]; exists {
		countQuery += fmt.Sprintf(" AND c.id = $%d", countArgCounter)
		countArgs = append(countArgs, class)
		countArgCounter++
	}
	if submissionType, exists := filter["type"]; exists {
		if submissionType == "this_week" {
			countQuery += " AND s.submitted_at >= EXTRACT(EPOCH FROM NOW() - INTERVAL '7 days')"
		} else if submissionType == "old" {
			countQuery += " AND s.submitted_at < EXTRACT(EPOCH FROM NOW() - INTERVAL '7 days')"
		}
	}

	var total int64
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		log.Error("[Repo][ListAdmin] Error Count Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch quiz submissions count")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Error("[Repo][ListAdmin] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch quiz submissions")
	}
	defer rows.Close()

	var submissions []quizEntity.ListQuizSubmissionAdmin
	for rows.Next() {
		var submission quizEntity.ListQuizSubmissionAdmin
		err := rows.Scan(
			&submission.ID, &submission.SetID, &submission.Correct,
			&submission.Grade, &submission.SubmittedAt,
			&submission.Name, &submission.Lesson, &submission.Class,
		)
		if err != nil {
			log.Error("[Repo][ListAdmin] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan quiz submissions")
		}
		submission.SubmittedAt = helper.FormatUnixTime(submission.SubmittedAt)
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListAdmin] Error Rows: ", err)
		return nil, app.NewAppError(500, "error while iterating quiz submissions")
	}

	return &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  submissions,
	}, nil
}
