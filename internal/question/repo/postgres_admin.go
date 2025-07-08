package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/redis/go-redis/v9"
)

func (r *questionRepository) ListAdmin(ctx context.Context, filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	var questions []questionEntity.ListQuestionAdmin
	var total int64

	cacheKeyParts := []string{"cache:question:list"}
	for _, key := range []string{"is_quiz", "lesson", "class", "set", "lang", "search", "code"} {
		if val, ok := filter[key]; ok {
			cacheKeyParts = append(cacheKeyParts, fmt.Sprintf("%s:%s", key, val))
		}
	}
	cacheKeyParts = append(cacheKeyParts, fmt.Sprintf("page:%d", page), fmt.Sprintf("limit:%d", limit))
	redisKey := strings.Join(cacheKeyParts, "|")

	if r.redis != nil {
		cached, err := cache.Get(ctx, redisKey, r.redis)
		if err != nil && err != redis.Nil {
			log.Warn("[Repo][ListAdmin] Redis error:", err)
		}
		if err == nil && cached != "" {
			var cachedResp response.PaginateResponse
			if err := json.Unmarshal([]byte(cached), &cachedResp); err == nil {
				return &cachedResp, nil
			} else {
				log.Warn("[Repo][ListAdmin] Failed to unmarshal cache:", err)
			}
		}
	}

	// SQL query
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
		parsed, err := strconv.ParseBool(isQuiz)
		if err != nil {
			log.Warn("[Repo][ListAdmin] Invalid value for is_quiz:", isQuiz)
			return nil, app.NewAppError(400, "invalid value for is_quiz")
		}
		whereClause += fmt.Sprintf(" AND q.is_quiz = $%d", argCounter)
		args = append(args, parsed)
		argCounter++
	}

	for _, key := range []string{"lesson", "class", "set", "lang"} {
		if val, exists := filter[key]; exists {
			column := map[string]string{
				"lesson": "l.name",
				"class":  "c.name",
				"set":    "s.name",
				"lang":   "q.lang",
			}[key]
			whereClause += fmt.Sprintf(" AND %s = $%d", column, argCounter)
			args = append(args, val)
			argCounter++
		}
	}

	if search, exists := filter["search"]; exists {
		whereClause += fmt.Sprintf(" AND q.content ILIKE $%d", argCounter)
		args = append(args, "%"+search+"%")
		argCounter++
	}

	if code, exists := filter["code"]; exists {
		whereClause += fmt.Sprintf(" AND l.code ILIKE $%d", argCounter)
		args = append(args, "%"+code+"%")
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

	// Pagination
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

	// Query
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

	// Response
	resp := &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  questions,
	}

	// Set to Redis
	if r.redis != nil {
		if dataJSON, err := json.Marshal(resp); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis, cache.ExpSecond)
		} else {
			log.Warn("[Repo][ListAdmin] Failed to marshal cache:", err)
		}
	}

	return resp, nil
}
