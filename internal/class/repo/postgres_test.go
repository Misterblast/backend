package repo_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ghulammuzz/misterblast/internal/class/entity"
	"github.com/ghulammuzz/misterblast/internal/class/repo"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repository := repo.NewClassRepository(mockDB, nil)

	mock.ExpectQuery("SELECT 1 FROM classes WHERE name = \\$1").
		WithArgs("Math").
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	exists, err := repository.Exists("Math")
	assert.NoError(t, err)
	assert.True(t, exists)

	mock.ExpectQuery("SELECT 1 FROM classes WHERE name = \\$1").
		WithArgs("Science").
		WillReturnError(sql.ErrNoRows)

	exists, err = repository.Exists("Science")
	assert.NoError(t, err)
	assert.False(t, exists)

	mock.ExpectQuery("SELECT 1 FROM classes WHERE name = \\$1").
		WithArgs("History").
		WillReturnError(sql.ErrConnDone)

	exists, err = repository.Exists("History")
	assert.Error(t, err)
	assert.False(t, exists)
}

func TestAddClass(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repository := repo.NewClassRepository(mockDB, nil)

	validClass := entity.SetClass{Name: "1"}
	invalidClass := entity.SetClass{Name: "Invalid"}

	mock.ExpectExec("INSERT INTO classes").
		WithArgs(validClass.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repository.Add(validClass)
	assert.NoError(t, err)

	err = repository.Add(invalidClass)
	assert.Error(t, err)
	assert.Equal(t, "validation failed", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteClass(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repository := repo.NewClassRepository(mockDB, nil)

	mock.ExpectExec("DELETE FROM classes WHERE id = ").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repository.Delete(1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListClass(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	repository := repo.NewClassRepository(mockDB, nil)

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "1").
		AddRow(2, "2")

	mock.ExpectQuery("SELECT id, name FROM classes").WillReturnRows(rows)

	classes, err := repository.List(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, classes, 2)
	assert.Equal(t, "1", classes[0].Name)
	assert.Equal(t, "2", classes[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
