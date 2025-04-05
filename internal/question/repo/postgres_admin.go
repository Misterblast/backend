package repo

import (
	"fmt"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

func (r *questionRepository) ListAdmin(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	baseQuery := `
		FROM questions q
		JOIN sets s ON q.set_id = s.id
		JOIN lessons l ON s.lesson_id = l.id
		JOIN classes c ON s.class_id = c.id
		WHERE 1=1 and deleted_at IS NULL
	`

	whereClause := ""
	args := []interface{}{}
	argCounter := 1

	if isQuiz, exists := filter["is_quiz"]; exists {
		whereClause += fmt.Sprintf(" AND q.is_quiz = $%d", argCounter)
		args = append(args, isQuiz)
		argCounter++
	}
	if lesson, exists := filter["lesson"]; exists {
		whereClause += fmt.Sprintf(" AND l.name = $%d", argCounter)
		args = append(args, lesson)
		argCounter++
	}
	if class, exists := filter["class"]; exists {
		whereClause += fmt.Sprintf(" AND c.name = $%d", argCounter)
		args = append(args, class)
		argCounter++
	}
	if set, exists := filter["set"]; exists {
		whereClause += fmt.Sprintf(" AND s.name = $%d", argCounter)
		args = append(args, set)
		argCounter++
	}

	// Count total
	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var total int64
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Error("[Repo][ListAdmin] Error Count Query:", err)
		return nil, app.NewAppError(500, "failed to count admin questions")
	}

	// Query data with pagination
	query := `
		SELECT q.id, q.number, q.type, q.format, q.content, q.explanation, q.is_quiz, q.set_id,
			   s.name AS set_name, l.name AS lesson_name, c.name AS class_name
	` + baseQuery + whereClause + " ORDER BY q.number"

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
		log.Error("[Repo][ListAdmin] Error Query:", err)
		return nil, app.NewAppError(500, "failed to fetch admin questions")
	}
	defer rows.Close()

	var questions []questionEntity.ListQuestionAdmin
	for rows.Next() {
		var q questionEntity.ListQuestionAdmin
		err := rows.Scan(&q.ID, &q.Number, &q.Type, &q.Format, &q.Content, &q.Explanation, &q.IsQuiz, &q.SetID, &q.SetName, &q.LessonName, &q.ClassName)
		if err != nil {
			log.Error("[Repo][ListAdmin] Error Scan:", err)
			return nil, app.NewAppError(500, "failed to scan admin questions")
		}
		questions = append(questions, q)
	}

	response := &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  questions,
	}

	return response, nil
}
