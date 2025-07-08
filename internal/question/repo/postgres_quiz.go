package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/redis/go-redis/v9"
)

func (r *questionRepository) AddQuizAnswer(answer questionEntity.SetAnswer) error {
	query := `
		INSERT INTO answers (question_id, code, content, img_url, is_answer) 
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, answer.QuestionID, answer.Code, answer.Content, answer.ImgURL, answer.IsAnswer)
	if err != nil {
		log.Error("[Repo][AddQuizAnswer] Error inserting answer: ", err)
		return app.NewAppError(500, "failed to insert quiz answer")
	}

	return nil
}

func (r *questionRepository) ListQuizQuestions(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionQuiz, error) {
	redisKeyParts := []string{"cache:quiz:list"}
	for _, key := range []string{"set_id", "type", "number"} {
		if val, ok := filter[key]; ok {
			redisKeyParts = append(redisKeyParts, fmt.Sprintf("%s=%s", key, val))
		}
	}
	redisKey := strings.Join(redisKeyParts, "|")

	if r.redis != nil {
		cached, err := cache.Get(ctx, redisKey, r.redis)
		if err != nil && err != redis.Nil {
			log.Warn("[Repo][ListQuizQuestions] Redis error:", err)
		}
		if err == nil && cached != "" {
			var cachedData []questionEntity.ListQuestionQuiz
			if err := json.Unmarshal([]byte(cached), &cachedData); err == nil {
				return cachedData, nil
			} else {
				log.Warn("[Repo][ListQuizQuestions] Unmarshal cache error:", err)
			}
		}
	}

	query := `
		SELECT q.id, q.number, q.type, q.format, q.content, q.set_id,
			   COALESCE(a.id, 0), COALESCE(a.code, ''), 
			   COALESCE(a.content, ''), COALESCE(a.img_url, '')
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id
		WHERE q.is_quiz = true
	`
	args := []interface{}{}
	argCounter := 1

	if setID, exists := filter["set_id"]; exists {
		query += fmt.Sprintf(" AND q.set_id = $%d", argCounter)
		args = append(args, setID)
		argCounter++
	}
	if questionType, exists := filter["type"]; exists {
		query += fmt.Sprintf(" AND q.type = $%d", argCounter)
		args = append(args, questionType)
		argCounter++
	}
	if number, exists := filter["number"]; exists {
		query += fmt.Sprintf(" AND q.number = $%d", argCounter)
		args = append(args, number)
		argCounter++
	}

	query += " ORDER BY q.number, a.code"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Error("[Repo][ListQuizQuestions] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch quiz questions")
	}
	defer rows.Close()

	questionsMap := make(map[int32]*questionEntity.ListQuestionQuiz)
	var questions []*questionEntity.ListQuestionQuiz

	for rows.Next() {
		var qID int32
		var number int
		var qType, qFormat, content string
		var setID int32
		var aID int32
		var code, aContent, imgURL string

		err := rows.Scan(&qID, &number, &qType, &qFormat, &content, &setID, &aID, &code, &aContent, &imgURL)
		if err != nil {
			log.Error("[Repo][ListQuizQuestions] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan quiz questions")
		}

		if _, exists := questionsMap[qID]; !exists {
			questionsMap[qID] = &questionEntity.ListQuestionQuiz{
				ID:      qID,
				Number:  number,
				Type:    qType,
				Format:  qFormat,
				Content: content,
				SetID:   setID,
				Answers: []questionEntity.ListAnswer{},
			}
			questions = append(questions, questionsMap[qID])
		}

		if aID != 0 {
			answer := questionEntity.ListAnswer{
				ID:      aID,
				Code:    code,
				Content: aContent,
			}
			if imgURL != "" {
				answer.ImgURL = &imgURL
			}
			questionsMap[qID].Answers = append(questionsMap[qID].Answers, answer)
		}
	}

	finalQuestions := make([]questionEntity.ListQuestionQuiz, len(questions))
	for i, q := range questions {
		finalQuestions[i] = *q
	}

	if r.redis != nil {
		if dataJSON, err := json.Marshal(finalQuestions); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis, cache.ExpBlazing)
		} else {
			log.Warn("[Repo][ListQuizQuestions] Failed to marshal for cache:", err)
		}
	}

	return finalQuestions, nil
}

func (r *questionRepository) DeleteAnswer(id int32) error {
	query := `DELETE FROM answers WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteAnswer] Error deleting answer:", err)
		return app.NewAppError(500, "failed to delete answer")
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return app.NewAppError(404, "answer not found")
	}

	return nil
}

func (r *questionRepository) EditAnswer(id int32, answer questionEntity.EditAnswer) error {
	query := `
		UPDATE answers 
		SET code = $1, content = $2, img_url = $3, is_answer = $4 
		WHERE id = $5`

	_, err := r.db.Exec(query, answer.Code, answer.Content, answer.ImgURL, answer.IsAnswer, id)
	if err != nil {
		log.Error("[Repo][EditAnswer] Error updating answer:", err)
		return app.NewAppError(500, err.Error())
	}

	return nil
}

func (r *questionRepository) ListQuizQuestionsLessonClass(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionQuiz, int, error) {
	var setID string

	if val, ok := filter["set_id"]; ok && val != "" {
		setID = val
	} else {
		lessonID, hasLesson := filter["lesson_id"]

		if !hasLesson {
			return nil, 0, app.NewAppError(400, "lesson_id is required if set_id is not provided")
		}

		var classID string
		queryClass := `
			SELECT class_id FROM (
				SELECT class_id FROM sets 
				WHERE is_quiz = true AND lesson_id = $1 
				GROUP BY class_id
				ORDER BY RANDOM()
				LIMIT 1
			) AS random_class
		`

		err := r.db.QueryRowContext(ctx, queryClass, lessonID).Scan(&classID)
		if err != nil {
			log.Error("[Repo][ListQuizQuestions] Failed to get random class_id: ", err)
			return nil, 0, app.NewAppError(404, "no class found for specified lesson")
		}

		querySet := `
			SELECT id FROM sets
			WHERE is_quiz = true AND lesson_id = $1 AND class_id = $2
			ORDER BY RANDOM()
			LIMIT 1
		`
		err = r.db.QueryRowContext(ctx, querySet, lessonID, classID).Scan(&setID)
		if err != nil {
			log.Error("[Repo][ListQuizQuestions] Failed to get random set_id: ", err)
			return nil, 0, app.NewAppError(404, "no quiz set found for specified lesson and class")
		}
	}

	log.Debug("[Repo][ListQuizQuestions] Using set_id: ", setID)

	redisKey := fmt.Sprintf("quiz:list:%s:%s:%s", setID, filter["type"], filter["number"])
	if r.redis != nil {
		if cached, err := cache.Get(ctx, redisKey, r.redis); err == nil && cached != "" {
			var cachedData []questionEntity.ListQuestionQuiz
			if err := json.Unmarshal([]byte(cached), &cachedData); err == nil {
				return cachedData, 0, nil
			}
		}
	}

	query := `
		SELECT q.id, q.number, q.type, q.format, q.content, q.set_id,
			   COALESCE(a.id, 0) AS answer_id, COALESCE(a.code, '') AS code, 
			   COALESCE(a.content, '') AS answer_content, COALESCE(a.img_url, '') AS img_url
		FROM questions q
		LEFT JOIN answers a ON q.id = a.question_id
		WHERE q.is_quiz = true AND q.set_id = $1
	`
	args := []interface{}{setID}
	argCounter := 2

	if questionType, exists := filter["type"]; exists && questionType != "" {
		query += fmt.Sprintf(" AND q.type = $%d", argCounter)
		args = append(args, questionType)
		argCounter++
	}
	if number, exists := filter["number"]; exists && number != "" {
		query += fmt.Sprintf(" AND q.number = $%d", argCounter)
		args = append(args, number)
		argCounter++
	}
	if lang, exists := filter["lang"]; exists && lang != "" {
		query += fmt.Sprintf(" AND q.lang = $%d", argCounter)
		args = append(args, lang)
		argCounter++
	}

	query += " ORDER BY q.id, a.code"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Error("[Repo][ListQuizQuestions] Error Query: ", err)
		return nil, 0, app.NewAppError(500, "failed to fetch quiz questions")
	}
	defer rows.Close()

	questionsMap := make(map[int32]*questionEntity.ListQuestionQuiz)
	var questions []*questionEntity.ListQuestionQuiz

	for rows.Next() {
		var qID, aID, setIDInt int32
		var number int
		var qType, qFormat, content, code, aContent, imgURL string

		err := rows.Scan(&qID, &number, &qType, &qFormat, &content, &setIDInt, &aID, &code, &aContent, &imgURL)
		if err != nil {
			log.Error("[Repo][ListQuizQuestions] Error Scan: ", err)
			return nil, 0, app.NewAppError(500, "failed to scan quiz questions")
		}

		if _, exists := questionsMap[qID]; !exists {
			questionsMap[qID] = &questionEntity.ListQuestionQuiz{
				ID:      qID,
				Number:  number,
				Type:    qType,
				Format:  qFormat,
				Content: content,
				SetID:   setIDInt,
				Answers: []questionEntity.ListAnswer{},
			}
			questions = append(questions, questionsMap[qID])
		}

		if aID != 0 {
			answer := questionEntity.ListAnswer{
				ID:      aID,
				Code:    code,
				Content: aContent,
			}
			if imgURL != "" {
				answer.ImgURL = &imgURL
			}
			questionsMap[qID].Answers = append(questionsMap[qID].Answers, answer)
		}
	}

	finalQuestions := make([]questionEntity.ListQuestionQuiz, len(questions))
	for i, q := range questions {
		finalQuestions[i] = *q
	}

	rand.Shuffle(len(finalQuestions), func(i, j int) {
		finalQuestions[i], finalQuestions[j] = finalQuestions[j], finalQuestions[i]
	})

	if r.redis != nil {
		if dataJSON, err := json.Marshal(finalQuestions); err == nil {
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis, cache.ExpBlazing)
		}
	}

	// convert string to int setID
	setIDInt, err := strconv.Atoi(setID)
	if err != nil {
		log.Error("[Repo][ListQuizQuestions] Error converting setID to int: ", err)
		return nil, 0, app.NewAppError(500, "failed to convert set_id to integer")
	}

	return finalQuestions, setIDInt, nil
}
