package repo

import (
	"fmt"

	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

func (r *quizRepository) ListAdmin(filter map[string]string, page, limit int) ([]quizEntity.ListQuizSubmissionAdmin, error) {
	query := `
		SELECT s.id, s.set_id, s.correct, s.grade, s.submitted_at,
			   u.name AS user_name,
			   l.name AS lesson_name, c.name AS class_name
		FROM quiz_submissions s
		JOIN users u ON s.user_id = u.id
		JOIN sets a ON s.set_id = a.id
		JOIN lessons l ON a.lesson_id = l.id
		JOIN classes c ON a.class_id = c.id
	`

	args := []interface{}{}
	argCounter := 1

	if lesson, exists := filter["lesson"]; exists {
		query += fmt.Sprintf(" AND l.name = $%d", argCounter)
		args = append(args, lesson)
		argCounter++
	}
	if class, exists := filter["class"]; exists {
		query += fmt.Sprintf(" AND c.name = $%d", argCounter)
		args = append(args, class)
		argCounter++
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

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Error("[Repo][ListAdmin] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch quiz submissions")
	}
	defer rows.Close()

	var submissions []quizEntity.ListQuizSubmissionAdmin
	for rows.Next() {
		var submission quizEntity.ListQuizSubmissionAdmin
		err := rows.Scan(&submission.ID, &submission.SetID, &submission.Correct, &submission.Grade, &submission.SubmittedAt, &submission.Name, &submission.Lesson, &submission.Class)
		if err != nil {
			log.Error("[Repo][ListAdmin] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan quiz submissions")
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListAdmin] Error Rows: ", err)
		return nil, app.NewAppError(500, "error while iterating quiz submissions")
	}

	return submissions, nil
}
