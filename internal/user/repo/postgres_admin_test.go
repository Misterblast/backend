package repo_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"

)

func TestUserRepository_AdminActivation(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := userRepo.NewUserRepository(mockDB)
	adminID := int32(1)

	mock.ExpectExec("UPDATE users SET is_verified=true WHERE id=\\$1").
		WithArgs(adminID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.AdminActivation(adminID)
	assert.NoError(t, err)
}

func TestUserRepository_Auth(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := userRepo.NewUserRepository(mockDB)
	id := int32(1)

	mock.ExpectQuery("SELECT id, name, email, COALESCE\\(img_url, ''\\), is_admin, is_verified FROM users WHERE id=\\$1").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "img_url", "is_admin", "is_verified"}).
			AddRow(1, "John Doe", "john@example.com", "", false, true))

	user, err := repo.Auth(id)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
}
