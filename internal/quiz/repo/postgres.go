package repo

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"

	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type QuizRepository interface {
	Submit(req quizEntity.QuizSubmit, setId int, userId int) error
	List(filter map[string]string, page, limit, userID int) ([]quizEntity.ListQuizSubmission, error)
	ListAdmin(filter map[string]string, page, limit int) ([]quizEntity.ListQuizSubmissionAdmin, error)
	GetLast(userID int) (quizEntity.QuizExp, error)
}

func (r *quizRepository) List(filter map[string]string, page, limit, userID int) ([]quizEntity.ListQuizSubmission, error) {
	query := `
		SELECT s.id, s.set_id, s.correct, s.grade, s.submitted_at,
			   l.name AS lesson_name, c.name AS class_name
		FROM quiz_submissions s
		JOIN sets a ON s.set_id = a.id
		JOIN lessons l ON a.lesson_id = l.id
		JOIN classes c ON a.class_id = c.id  
		WHERE s.user_id = $1
	`

	args := []interface{}{userID}
	argCounter := 2

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
		log.Error("[Repo][List] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch quiz submissions")
	}
	defer rows.Close()

	var submissions []quizEntity.ListQuizSubmission
	for rows.Next() {
		var submission quizEntity.ListQuizSubmission
		err := rows.Scan(&submission.ID, &submission.SetID, &submission.Correct, &submission.Grade, &submission.SubmittedAt, &submission.Lesson, &submission.Class)
		if err != nil {
			log.Error("[Repo][List] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan quiz submissions")
		}
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][List] Error Rows: ", err)
		return nil, app.NewAppError(500, "error while iterating quiz submissions")
	}

	return submissions, nil
}

func (r *quizRepository) checkTotalQuestion(setID int) (int, error) {
	var total int

	query := `SELECT COUNT(*) FROM questions WHERE set_id = $1`

	err := r.db.QueryRow(query, setID).Scan(&total)
	if err != nil {
		log.Error("[Repo][checkTotalQuestion] Error Exec: ", err)
		return 0, app.NewAppError(500, err.Error())
	}

	return total, nil
}

func (r *quizRepository) checkCorrectAnswer(setID int) (string, error) {
	var correctAnswers string

	query := `
	SELECT STRING_AGG(a.code, '' ORDER BY q.number) AS correct_answers
		FROM questions q
		JOIN answers a ON q.id = a.question_id
		WHERE q.set_id = $1 AND a.is_answer = true
		`

	err := r.db.QueryRow(query, setID).Scan(&correctAnswers)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		log.Error("[Repo][checkCorrectAnswer] Error Exec: ", err)
		return "", app.NewAppError(500, err.Error())
	}

	return correctAnswers, nil
}

func (r *quizRepository) checkQuizScore(userAnswer, correctAnswer string, setID int) (int, int, error) {
	totalQuestions, err := r.checkTotalQuestion(setID)
	if err != nil {
		return 0, 0, err
	}

	if totalQuestions == 0 {
		log.Error("[Repo][checkQuizScore] No questions found in this set")
		return 0, 0, app.NewAppError(400, "no questions found in this set")
	}

	correctCount := 0
	for i := 0; i < len(userAnswer) && i < len(correctAnswer); i++ {
		if userAnswer[i] == correctAnswer[i] {
			correctCount++
		}
	}

	score := (correctCount * 100) / totalQuestions

	return score, correctCount, nil
}

func (r *quizRepository) getNextAttemptNo(setID int, userID int) (int, error) {
	var attemptNo int
	err := r.db.QueryRow("SELECT COALESCE(MAX(attempt_no), 0) + 1 FROM quiz_submissions WHERE user_id = $1 AND set_id = $2", userID, setID).Scan(&attemptNo)
	if err != nil {
		log.Error("[Repo][getNextAttemptNo] Error Exec: ", err)
		return 0, app.NewAppError(500, err.Error())
	}
	return attemptNo, nil
}

func (r *quizRepository) Submit(req quizEntity.QuizSubmit, setID int, userID int) error {
	correctAnswer, err := r.checkCorrectAnswer(setID)
	if err != nil {
		return err
	}

	totalQuestions, err := r.checkTotalQuestion(setID)
	if err != nil {
		return err
	}
	if len(req.Answers) != totalQuestions {
		log.Error("[Repo][Submit] Invalid number of answers provided")
		return app.NewAppError(400, "invalid number of answers provided")
	}

	sort.Slice(req.Answers, func(i, j int) bool {
		return req.Answers[i].Number < req.Answers[j].Number
	})

	var answers []string
	for _, ans := range req.Answers {
		answers = append(answers, ans.Answer)
	}
	answerStr := strings.Join(answers, "")

	attemptNo, err := r.getNextAttemptNo(setID, userID)
	if err != nil {
		return err
	}

	score, correctCount, err := r.checkQuizScore(answerStr, correctAnswer, setID)
	if err != nil {
		return err
	}

	query := "INSERT INTO quiz_submissions (answer, correct, grade, attempt_no, set_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err = r.db.Exec(query, answerStr, correctCount, score, attemptNo, setID, userID)
	if err != nil {
		log.Error("[Repo][Submit] Error Exec: ", err)
		return app.NewAppError(500, err.Error())
	}
	return nil
}

type quizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) QuizRepository {
	return &quizRepository{db: db}
}
