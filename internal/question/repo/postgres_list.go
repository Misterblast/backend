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
)

func (r *questionRepository) List(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	var questions []questionEntity.ListQuestionExample

	cacheKeyParts := []string{"question-user:list"}
	for _, key := range []string{"is_quiz", "lesson", "class", "set", "lang", "search"} {
		cacheKeyParts = append(cacheKeyParts, fmt.Sprintf("%s:%s", key, filter[key]))
	}
	redisKey := strings.Join(cacheKeyParts, "|")

	if r.redis != nil {
		cached, err := cache.Get(ctx, redisKey, r.redis)
		if err == nil && cached != "" {
			var cachedQuestions []questionEntity.ListQuestionExample
			if err := json.Unmarshal([]byte(cached), &cachedQuestions); err == nil {
				return cachedQuestions, nil
			}
		}
	}

	baseQuery := `
		FROM questions q
		JOIN sets s ON q.set_id = s.id
		JOIN lessons l ON s.lesson_id = l.id
		JOIN classes c ON s.class_id = c.id
		LEFT JOIN answers a ON q.id = a.question_id
		WHERE q.deleted_at IS NULL
	`

	whereClause := ""
	args := []interface{}{}
	argCounter := 1

	if isQuiz, exists := filter["is_quiz"]; exists {
		whereClause += fmt.Sprintf(" AND q.is_quiz = $%d", argCounter)
		args = append(args, isQuiz)
		argCounter++
	}
	if lesson, exists := filter["lesson_id"]; exists {
		whereClause += fmt.Sprintf(" AND l.id = $%d", argCounter)
		args = append(args, lesson)
		argCounter++
	}
	if class, exists := filter["class_id"]; exists {
		whereClause += fmt.Sprintf(" AND c.id = $%d", argCounter)
		args = append(args, class)
		argCounter++
	}
	if set, exists := filter["set_id"]; exists {
		whereClause += fmt.Sprintf(" AND s.id = $%d", argCounter)
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

	query := `
		SELECT 
			q.id, q.number, q.type, q.format, q.content, q.explanation, q.reasoning, q.set_id,
			COALESCE(
				json_agg(
					json_build_object(
						'id', a.id,
						'code', a.code,
						'content', a.content,
						'img_url', a.img_url
					)
				) FILTER (WHERE a.id IS NOT NULL), '[]'
			) AS answers
	` + baseQuery + whereClause + `
		GROUP BY q.id
		ORDER BY q.number
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Error("[Repo][List] Error Query:", err)
		return nil, app.NewAppError(500, "failed to fetch questions")
	}
	defer rows.Close()

	for rows.Next() {
		var q questionEntity.ListQuestionExample
		var answersJSON []byte

		err := rows.Scan(
			&q.ID,
			&q.Number,
			&q.Type,
			&q.Format,
			&q.Content,
			&q.Explanation,
			&q.Reason,
			&q.SetID,
			&answersJSON,
		)
		if err != nil {
			log.Error("[Repo][List] Error Scan:", err)
			return nil, app.NewAppError(500, "failed to scan question")
		}

		if err := json.Unmarshal(answersJSON, &q.Answers); err != nil {
			log.Error("[Repo][List] Error Unmarshal Answers:", err)
			return nil, app.NewAppError(500, "failed to parse answers")
		}

		questions = append(questions, q)
	}

	if r.redis != nil {
		if dataJSON, err := json.Marshal(questions); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis, cache.ExpSecond)
		}
	}

	return questions, nil
}
