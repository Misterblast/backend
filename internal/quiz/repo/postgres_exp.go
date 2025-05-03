package repo

import (
	"database/sql"
	"errors"

	quizEntity "github.com/ghulammuzz/misterblast/internal/quiz/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"

	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/lib/pq"
)

func (r *quizRepository) GetLast(userID int) (quizEntity.QuizExp, error) {
	var quiz quizEntity.QuizExp

	query := `
		SELECT id, answer, correct, grade, attempt_no, submitted_at, set_id
		FROM quiz_submissions
		WHERE user_id = $1
		ORDER BY submitted_at DESC
		LIMIT 1
	`
	var correct, attemptNo, setID int
	var answer string
	if err := r.db.QueryRow(query, userID).Scan(&quiz.ID, &answer, &correct, &quiz.Grade, &attemptNo, &quiz.SubmittedAt, &setID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quizEntity.QuizExp{}, nil
		}
		log.Error("[quizRepo.GetLast] failed to get last quiz submission", err.Error())
		return quizEntity.QuizExp{}, app.NewAppError(500, "failed to get last quiz submission")
	}

	var totalQuestions int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM questions WHERE set_id = $1", setID).Scan(&totalQuestions); err != nil {
		log.Error("[quizRepo.GetLast] failed to get total questions", err.Error())
		return quizEntity.QuizExp{}, app.NewAppError(500, "failed to get total questions")
	}

	// log.Debug("total questions: ", totalQuestions)
	// log.Debug("correct: ", correct)
	quiz.Correct = correct
	quiz.Wrong = totalQuestions - correct
	// log.Debug("wrong: ", quiz.Wrong)
	quiz.AttemptNo = attemptNo

	questionsQuery := `
		SELECT id, number, content, format, explanation
		FROM questions
		WHERE set_id = $1
		ORDER BY number ASC
	`
	rows, err := r.db.Query(questionsQuery, setID)
	if err != nil {
		log.Error("[quizRepo.GetLast] failed to get questions", err.Error())
		return quizEntity.QuizExp{}, app.NewAppError(500, "failed to get questions")
	}
	defer rows.Close()

	var questions []quizEntity.QuizExpObj
	questionMap := make(map[int]*quizEntity.QuizExpObj)
	questionIDs := []int{}

	for rows.Next() {
		var q quizEntity.QuizExpObj
		var questionID int
		if err := rows.Scan(&questionID, &q.Number, &q.QuestionContent, &q.Explanation, &q.Explanation); err != nil {
			log.Error("[quizRepo.GetLast] failed to scan questions", err.Error())
			return quizEntity.QuizExp{}, app.NewAppError(500, "failed to scan questions")
		}
		questionMap[questionID] = &q
		questionIDs = append(questionIDs, questionID)
	}

	answersQuery := `
		SELECT question_id, code, content
		FROM answers
		WHERE question_id = ANY($1)
		AND is_answer = TRUE
	`
	ansRows, err := r.db.Query(answersQuery, pq.Array(questionIDs))
	if err != nil {
		log.Error("[quizRepo.GetLast] failed to get answers", err.Error())
		return quizEntity.QuizExp{}, app.NewAppError(500, "failed to get answers")
	}
	defer ansRows.Close()

	actualAnswers := make(map[int]quizEntity.QuizExpObj)
	for ansRows.Next() {
		var questionID int
		var code, content string
		if err := ansRows.Scan(&questionID, &code, &content); err != nil {
			log.Error("[quizRepo.GetLast] failed to scan answers", err.Error())
			return quizEntity.QuizExp{}, app.NewAppError(500, "failed to scan answers")
		}
		actualAnswers[questionID] = quizEntity.QuizExpObj{
			ActualCode:    code,
			ActualContent: content,
		}
	}

	userAnswers := []rune(answer)
	for i, questionID := range questionIDs {
		if i >= len(userAnswers) {
			break
		}
		q := questionMap[questionID]
		actualAns := actualAnswers[questionID]
		q.UserCode = string(userAnswers[i])
		q.ActualCode = actualAns.ActualCode

		var userContent string
		userAnswerQuery := `
			SELECT content FROM answers 
			WHERE question_id = $1 AND code = $2
		`
		if err := r.db.QueryRow(userAnswerQuery, questionID, q.UserCode).Scan(&userContent); err != nil {
			userContent = ""
		}
		q.UserContent = userContent
		q.ActualContent = actualAns.ActualContent
		q.IsCorrect = q.UserCode == q.ActualCode
		questions = append(questions, *q)
	}

	quiz.Answers = questions
	return quiz, nil
}
