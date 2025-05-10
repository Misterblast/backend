package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

func (r *questionRepository) ListAdmin(ctx context.Context, filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	var questions []questionEntity.ListQuestionAdmin
	var total int64

	cacheKeyParts := []string{"question:list"}
	for _, key := range []string{"is_quiz", "lesson", "class", "set", "lang", "search"} {
		cacheKeyParts = append(cacheKeyParts, fmt.Sprintf("%s:%s", key, filter[key]))
	}
	cacheKeyParts = append(cacheKeyParts, fmt.Sprintf("page:%d", page), fmt.Sprintf("limit:%d", limit))
	redisKey := strings.Join(cacheKeyParts, "|")
	if r.redis != nil {
		cached, err := cache.Get(ctx, redisKey, r.redis)
		if err == nil && cached != "" {
			var cachedResp response.PaginateResponse
			if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
				return &cachedResp, nil
			}
		}
	}

	// cached, err := cache.Get(ctx, redisKey, r.redis)
	// if err == nil && cached != "" {
	// 	if err := json.Unmarshal([]byte(cached), &questions); err == nil {
	// 		return questions, nil
	// 	}
	// }

	baseQuery := `
		FROM questions q
		JOIN sets s ON q.set_id = s.id
		JOIN lessons l ON s.lesson_id = l.id
		JOIN classes c ON s.class_id = c.id
		WHERE 1=1 AND q.deleted_at IS NULL
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
	if lang, exists := filter["lang"]; exists {
		whereClause += fmt.Sprintf(" AND q.lang = $%d", argCounter)
		args = append(args, lang)
		argCounter++
	}
	if search, exists := filter["search"]; exists {
		whereClause += fmt.Sprintf(" AND q.content ILIKE $%d", argCounter)
		args = append(args, "%"+search+"%")
		argCounter++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Error("[Repo][ListAdmin] Error Count Query:", err)
		return nil, app.NewAppError(500, "failed to count admin questions")
	}

	query := `
		SELECT q.id, q.number, q.type, q.format, q.content, q.explanation, q.reasoning, q.is_quiz, q.set_id,
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

	for rows.Next() {
		var q questionEntity.ListQuestionAdmin
		err := rows.Scan(&q.ID, &q.Number, &q.Type, &q.Format, &q.Content, &q.Explanation, &q.Reason, &q.IsQuiz, &q.SetID, &q.SetName, &q.LessonName, &q.ClassName)
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

	if r.redis != nil {
		if dataJSON, err := json.Marshal(response); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis, cache.ExpBlazing)
		}
	}

	return response, nil
}
