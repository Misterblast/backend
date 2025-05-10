package repo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/repo"
	"github.com/stretchr/testify/assert"
)

func TestAddQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewQuestionRepository(db, nil)

	mock.ExpectExec(`INSERT INTO questions`).
		WithArgs(1, "c4_faktual", "mm", "Sample Question", 1, "exp-1", "id").
		WillReturnResult(sqlmock.NewResult(0, 1))

	question := questionEntity.SetQuestion{
		SetID:       1,
		Number:      1,
		Type:        "c4_faktual",
		Format:      "mm",
		Content:     "Sample Question",
		Explanation: "exp-1",
	}

	err = repository.Add(question, "id")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestDeleteQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewQuestionRepository(db, nil)

	const query = `UPDATE questions SET deleted_at = EXTRACT\(EPOCH FROM NOW\(\)\) WHERE id = \$1 AND deleted_at IS NULL`

	t.Run("Success delete", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repository.Delete(1)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("No rows affected", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(2).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repository.Delete(2)
		assert.Error(t, err)
		assert.Equal(t, "question not found or already deleted", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectExec(query).
			WithArgs(3).
			WillReturnError(errors.New("db error"))

		err := repository.Delete(3)
		assert.Error(t, err)
		assert.Equal(t, "failed to delete question", err.Error())
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestExistsQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewQuestionRepository(db, nil)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM questions WHERE set_id =`).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	exists, err := repository.Exists(1, 1)

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListAdmin(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewQuestionRepository(db, nil)

	// Mock total count query
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM questions`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock data query
	mockRows := sqlmock.NewRows([]string{"id", "number", "type", "format", "content", "explanation", "reason", "is_quiz", "set_id", "set_name", "lesson_name", "class_name"}).
		AddRow(1, 1, "c4_faktual", "mm", "Question 1", "exp-1", "r-1", true, 1, "Set 1", "Lesson 1", "Class 1")

	mock.ExpectQuery(`SELECT q.id, q.number, q.type, q.format, q.content, q.explanation, q.reasoning, q.is_quiz, q.set_id`).
		WillReturnRows(mockRows)

	result, err := repository.ListAdmin(context.Background(), map[string]string{}, 1, 10)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.Len(t, result.Data, 1)

	questions, ok := result.Data.([]questionEntity.ListQuestionAdmin)
	assert.True(t, ok)
	assert.Equal(t, "Question 1", questions[0].Content)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEditQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewQuestionRepository(db, nil)

	editQuestion := questionEntity.EditQuestion{SetID: 9, Number: 2, Type: "c4_faktual", Format: "mm", Content: "Updated Content", IsQuiz: false, Explanation: "exp-1"}

	mock.ExpectExec(`UPDATE questions SET number =`).
		WithArgs(editQuestion.Number, editQuestion.Type, editQuestion.Format, editQuestion.Content, editQuestion.IsQuiz, editQuestion.SetID, editQuestion.Explanation, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repository.Edit(1, editQuestion)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
