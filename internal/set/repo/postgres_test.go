package repo_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ghulammuzz/misterblast/internal/set/entity"
	"github.com/ghulammuzz/misterblast/internal/set/repo"
	"github.com/stretchr/testify/assert"
)

func TestAddSet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewSetRepository(db)

	mock.ExpectExec("INSERT INTO sets").
		WithArgs("Set A", 1, 1, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	set := entity.SetSet{Name: "Set A", LessonID: 1, ClassID: 1, IsQuiz: false}
	err = repository.Add(set)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSet(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewSetRepository(db)

	mock.ExpectExec(`DELETE FROM sets WHERE id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repository.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListSets(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewSetRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "lesson", "class", "is_quiz"}).
		AddRow(1, "Set A", "Math", "Class 1", false).
		AddRow(2, "Set B", "Science", "Class 2", true)

	mock.ExpectQuery("SELECT s.id, s.name, l.name AS lesson, c.name AS class, s.is_quiz FROM sets").
		WillReturnRows(rows)

	filter := map[string]string{}
	sets, err := repository.List(filter)
	assert.NoError(t, err)
	assert.Len(t, sets, 2)
	assert.Equal(t, "Set A", sets[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListWithFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repository := repo.NewSetRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "lesson", "class", "is_quiz"}).
		AddRow(1, "Set A", "Math", "Class 1", false)

	mock.ExpectQuery(`SELECT s.id, s.name, l.name AS lesson, c.name AS class, s.is_quiz FROM sets s`+
		` JOIN lessons l ON s.lesson_id = l.id`+
		` JOIN classes c ON s.class_id = c.id WHERE 1=1 AND l.name = \$1 AND c.name = \$2`).
		WithArgs("Math", "Class 1").
		WillReturnRows(rows)

	filters := map[string]string{"lesson": "Math", "class": "Class 1"}
	sets, err := repository.List(filters)
	assert.NoError(t, err)
	assert.Len(t, sets, 1)
	assert.Equal(t, "Set A", sets[0].Name)
}
