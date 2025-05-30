package repo

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ghulammuzz/misterblast/helper"
	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type QuizRepository interface {
	Submit(req quizEntity.QuizSubmit, setId int, userId int) (int, error)
	List(filter map[string]string, userID int) (*response.PaginateResponse, error)
	ListAdmin(filter map[string]string, page, limit int) (*response.PaginateResponse, error)
	GetLast(userID int) (quizEntity.QuizExp, error)
}

func (r *quizRepository) List(filter map[string]string, userID int) (*response.PaginateResponse, error) {
	page := 1
	limit := 10

	if val, ok := filter["page"]; ok {
		if p, err := strconv.Atoi(val); err == nil && p > 0 {
			page = p
		}
	}
	if val, ok := filter["limit"]; ok {
		if l, err := strconv.Atoi(val); err == nil && l > 0 {
			limit = l
		}
	}

	baseQuery := `
		FROM quiz_submissions s
		JOIN sets a ON s.set_id = a.id
		JOIN lessons l ON a.lesson_id = l.id
		JOIN classes c ON a.class_id = c.id
		WHERE s.user_id = $1
	`

	args := []interface{}{userID}
	argCounter := 2

	if lesson, exists := filter["lesson_id"]; exists {
		baseQuery += fmt.Sprintf(" AND l.id = $%d", argCounter)
		args = append(args, lesson)
		argCounter++
	}
	if class, exists := filter["class_id"]; exists {
		baseQuery += fmt.Sprintf(" AND c.id = $%d", argCounter)
		args = append(args, class)
		argCounter++
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Error("[Repo][List] Error Count Query: ", err)
		return nil, app.NewAppError(500, "failed to count quiz submissions")
	}

	mainQuery := `
		SELECT s.id, s.set_id, s.correct, s.grade, s.submitted_at,
			   l.name AS lesson_name, c.name AS class_name
	` + baseQuery + " ORDER BY s.submitted_at DESC"

	mainArgs := append([]interface{}{}, args...)

	if limit > 0 {
		mainQuery += fmt.Sprintf(" LIMIT $%d", argCounter)
		mainArgs = append(mainArgs, limit)
		argCounter++
	}
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		mainQuery += fmt.Sprintf(" OFFSET $%d", argCounter)
		mainArgs = append(mainArgs, offset)
	}

	rows, err := r.db.Query(mainQuery, mainArgs...)
	if err != nil {
		log.Error("[Repo][List] Error Query: ", err)
		return nil, app.NewAppError(500, "failed to fetch quiz submissions")
	}
	defer rows.Close()

	var submissions []quizEntity.ListQuizSubmission
	for rows.Next() {
		var submission quizEntity.ListQuizSubmission
		err := rows.Scan(
			&submission.ID, &submission.SetID, &submission.Correct,
			&submission.Grade, &submission.SubmittedAt,
			&submission.Lesson, &submission.Class,
		)
		if err != nil {
			log.Error("[Repo][List] Error Scan: ", err)
			return nil, app.NewAppError(500, "failed to scan quiz submissions")
		}
		submission.SubmittedAt = helper.FormatUnixTime(submission.SubmittedAt)
		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		log.Error("[Repo][List] Error Rows: ", err)
		return nil, app.NewAppError(500, "error while iterating quiz submissions")
	}

	paginateResp := &response.PaginateResponse{
		Total: total,
		Page:  page,
		Limit: limit,
		Data:  submissions,
	}

	return paginateResp, nil
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

func (r *quizRepository) Submit(req quizEntity.QuizSubmit, setID int, userID int) (int, error) {
	correctAnswer, err := r.checkCorrectAnswer(setID)
	if err != nil {
		return 0, err
	}

	totalQuestions, err := r.checkTotalQuestion(setID)
	if err != nil {
		return 0, err
	}
	if len(req.Answers) != totalQuestions {
		log.Error("[Repo][Submit] Invalid number of answers provided")
		return 0, app.NewAppError(400, "invalid number of answers provided")
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
		return 0, err
	}

	score, correctCount, err := r.checkQuizScore(answerStr, correctAnswer, setID)
	if err != nil {
		return 0, err
	}

	var id int
	query := "INSERT INTO quiz_submissions (answer, correct, grade, attempt_no, set_id, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	err = r.db.QueryRow(query, answerStr, correctCount, score, attemptNo, setID, userID).Scan(&id)
	if err != nil {
		log.Error("[Repo][Submit] Error Exec: ", err)
		return 0, app.NewAppError(500, err.Error())
	}
	return id, nil
}

type quizRepository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) QuizRepository {
	return &quizRepository{db: db}
}
