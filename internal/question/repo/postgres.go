package repo

import (
	"database/sql"
	"fmt"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type QuestionRepository interface {
	// Questions
	Add(question questionEntity.SetQuestion) error
	List(filter map[string]string) ([]questionEntity.ListQuestionExample, error)
	Delete(id int32) error
	Detail(id int32) (questionEntity.DetailQuestionExample, error)
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
	db *sql.DB
}

func NewQuestionRepository(db *sql.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Add(question questionEntity.SetQuestion) error {
	query := `INSERT INTO questions (number, type, format, content, is_quiz, explanation, set_id) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, question.Number, question.Type, question.Format, question.Content, question.IsQuiz, question.Explanation, question.SetID)
	if err != nil {
		log.Error("[Repo][AddQuestion] Error inserting question:", err)
		return app.NewAppError(500, err.Error())
	}
	return nil
}

func (r *questionRepository) Detail(id int32) (questionEntity.DetailQuestionExample, error) {
	query := `SELECT id, number, type, format, content, explanation, set_id FROM questions WHERE id = $1`
	var question questionEntity.DetailQuestionExample
	err := r.db.QueryRow(query, id).Scan(&question.ID, &question.Number, &question.Type, &question.Format, &question.Content, &question.Explanation, &question.SetID)
	if err != nil {
		if err == sql.ErrNoRows {
			return question, app.NewAppError(404, "question not found")
		}
		log.Error("[Repo][DetailQuestion] Error fetching question detail:", err)
		return question, app.NewAppError(500, "failed to fetch question detail")
	}
	return question, nil
}

func (r *questionRepository) List(filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	query := `SELECT id, number, type, format, content, explanation, set_id FROM questions WHERE 1=1`
	args := []interface{}{}
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

	var questions []questionEntity.ListQuestionExample
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

	return questions, nil
}

func (r *questionRepository) Delete(id int32) error {
	query := `DELETE FROM questions WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("[Repo][DeleteQuestion] Error deleting question:", err)
		return app.NewAppError(500, "failed to delete question")
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return app.NewAppError(404, "question not found")
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
