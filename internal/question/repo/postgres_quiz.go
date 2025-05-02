package repo

import (
	"context"
	"encoding/json"
	"fmt"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
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
	redisKey := fmt.Sprintf("quiz:list:%s:%s:%s",
		filter["set_id"], filter["type"], filter["number"],
	)
	if r.redis != nil {
		cached, err := cache.Get(ctx, redisKey, r.redis)
		if err == nil && cached != "" {
			var cachedData []questionEntity.ListQuestionQuiz
			if err := json.Unmarshal([]byte(cached), &cachedData); err == nil {
				return cachedData, nil
			}
		}
	}
	query := `
		SELECT q.id, q.number, q.type, q.format, q.content, q.set_id,
			   COALESCE(a.id, 0) AS answer_id, COALESCE(a.code, '') AS code, 
			   COALESCE(a.content, '') AS answer_content, COALESCE(a.img_url, '') AS img_url
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
		var code, aContent string
		var imgURL string

		err := rows.Scan(&qID, &number, &qType, &qFormat, &content, &setID,
			&aID, &code, &aContent, &imgURL)
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
			_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis)
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
		SET question_id = $1, code = $2, content = $3, img_url = $4, is_answer = $5 
		WHERE id = $6`

	_, err := r.db.Exec(query, answer.QuestionID, answer.Code, answer.Content, answer.ImgURL, answer.IsAnswer, id)
	if err != nil {
		log.Error("[Repo][EditAnswer] Error updating answer:", err)
		return app.NewAppError(500, err.Error())
	}

	return nil
}
