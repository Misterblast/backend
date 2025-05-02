package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	cache "github.com/ghulammuzz/misterblast/config/redis"
	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/redis/go-redis/v9"
)

type QuestionRepository interface {
	// Questions
	Add(question questionEntity.SetQuestion, lang string) error
	List(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionExample, error)
	Delete(id int32) error
	Detail(ctx context.Context, id int32) (questionEntity.DetailQuestionExample, error)
	Exists(setID int32, number int) (bool, error)
	Edit(id int32, question questionEntity.EditQuestion) error

	// Answer
	AddQuizAnswer(answer questionEntity.SetAnswer) error
	ListQuizQuestions(filter map[string]string) ([]questionEntity.ListQuestionQuiz, error)
	DeleteAnswer(id int32) error
	EditAnswer(id int32, answer questionEntity.EditAnswer) error

	// Admin
	ListAdmin(filter map[string]string, page, limit int) (*response.PaginateResponse, error)
}

type questionRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewQuestionRepository(db *sql.DB, redis *redis.Client) QuestionRepository {
	return &questionRepository{db, redis}
}

func (r *questionRepository) Add(question questionEntity.SetQuestion, lang string) error {
	query := `
		INSERT INTO questions (number, type, format, content, is_quiz, explanation, set_id, lang)
		VALUES ($1, $2, $3, $4, (SELECT is_quiz FROM sets WHERE id = $5), $6, $5, $7)
	`
	_, err := r.db.Exec(query,
		question.Number,
		question.Type,
		question.Format,
		question.Content,
		question.SetID,
		question.Explanation,
		lang,
	)

	if err != nil {
		log.Error("[Repo][AddQuestion] Error inserting question:", err)
		return app.NewAppError(500, err.Error())
	}
	return nil
}

func (r *questionRepository) Detail(ctx context.Context, id int32) (questionEntity.DetailQuestionExample, error) {
	var question questionEntity.DetailQuestionExample
	redisKey := fmt.Sprintf("question:detail:%d", id)

	cached, err := cache.Get(ctx, redisKey, r.redis)
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &question); err == nil {
			return question, nil
		}
	}

	var answersJSON []byte
	query := `
	SELECT 
		q.id, q.number, q.type, q.format, q.content, q.explanation, q.set_id,
		COALESCE(json_agg(json_build_object(
			'id', a.id,
			'code', a.code,
			'content', a.content,
			'img_url', a.img_url
		)) FILTER (WHERE a.id IS NOT NULL), '[]') AS answers
	FROM questions q
	LEFT JOIN answers a ON q.id = a.question_id
	WHERE q.id = $1
	GROUP BY q.id
	`

	err = r.db.QueryRow(query, id).Scan(
		&question.ID,
		&question.Number,
		&question.Type,
		&question.Format,
		&question.Content,
		&question.Explanation,
		&question.SetID,
		&answersJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return question, app.NewAppError(404, "question not found")
		}
		log.Error("[Repo][DetailQuestion] Error scanning:", err)
		return question, app.NewAppError(500, "failed to fetch question detail")
	}

	if err := json.Unmarshal(answersJSON, &question.Answers); err != nil {
		log.Error("[Repo][DetailQuestion] Error unmarshalling answers:", err)
		return question, app.NewAppError(500, "failed to parse answers")
	}

	dataJSON, err := json.Marshal(question)
	if err == nil {
		_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis)
		// _ = r.redis.Set(ctx, redisKey, dataJSON, 10*time.Minute).Err()
	}

	return question, nil
}

func (r *questionRepository) List(ctx context.Context,  filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	var questions []questionEntity.ListQuestionExample
	redisKey := fmt.Sprintf("question:list:%s", filter["set_id"])

	cached, err := cache.Get(ctx, redisKey, r.redis)
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &questions); err == nil {
			return questions, nil
		}
	}
	query := `SELECT id, number, type, format, content, explanation, set_id FROM questions WHERE 1=1 and deleted_at IS NULL`
	args := []any{}
	argCounter := 1

	if setID, ok := filter["set_id"]; ok {
		query += fmt.Sprintf(" AND set_id = $%d", argCounter)
		args = append(args, setID)
		argCounter++
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Error("[Repo][ListQuestions] Error Query: ", err)
		return nil, app.NewAppError(500, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var question questionEntity.ListQuestionExample
		if err := rows.Scan(&question.ID, &question.Number, &question.Type, &question.Format, &question.Content, &question.Explanation, &question.SetID); err != nil {
			log.Error("[Repo][ListQuestions] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan question")
		}
		questions = append(questions, question)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][ListQuestions] Error Iterating Rows: ", err)
		return nil, app.NewAppError(500, "error iterating rows")
	}
	dataJSON, err := json.Marshal(questions)
	if err == nil {
		_ = cache.Set(ctx, redisKey, string(dataJSON), r.redis)
		// _ = r.redis.Set(ctx, redisKey, dataJSON, 10*time.Minute).Err()
	}

	return questions, nil
}

func (r *questionRepository) Delete(id int32) error {
	const query = `
		UPDATE questions
		SET deleted_at = EXTRACT(EPOCH FROM NOW())
		WHERE id = $1 AND deleted_at IS NULL
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteQuestion] Failed to soft delete question with id %d: %v", id, err)
		return app.NewAppError(500, "failed to delete question")
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return app.NewAppError(404, "question not found or already deleted")
	}

	return nil
}

func (r *questionRepository) Exists(setID int32, number int) (bool, error) {
	query := `SELECT COUNT(*) FROM questions WHERE set_id = $1 AND number = $2`
	var count int
	err := r.db.QueryRow(query, setID, number).Scan(&count)
	if err != nil {
		log.Error("[Repo][ExistsQuestion] Error checking question:", err)
		return false, app.NewAppError(500, "failed to check question existence")
	}

	return count > 0, nil
}

func (r *questionRepository) Edit(id int32, question questionEntity.EditQuestion) error {
	query := `
		UPDATE questions 
		SET number = $1, type = $2, content = $4, format = $3, is_quiz = $5, set_id = $6, explanation = $7
		WHERE id = $8`

	_, err := r.db.Exec(query, question.Number, question.Type, question.Format, question.Content, question.IsQuiz, question.SetID, question.Explanation, id)
	if err != nil {
		log.Error("[Repo][EditQuestion] Error updating question:", err)
		return app.NewAppError(500, err.Error())
	}

	return nil
}
