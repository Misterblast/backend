package repo_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error initializing mock database: %v", err)
	}
	return mockDB, mock
}

func TestUserRepository_Exists(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := userRepo.NewUserRepository(mockDB)
	id := int32(1)

	mock.ExpectQuery("SELECT EXISTS\\(SELECT 1 FROM users WHERE id=\\$1\\)").
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.Exists(id)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestUserRepository_Add(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := userRepo.NewUserRepository(mockDB)
	user := userEntity.Register{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	isVerified := true

	mock.ExpectExec(`INSERT INTO users \(name, email, password, img_url, is_verified\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
		WithArgs(user.Name, user.Email, sqlmock.AnyArg(), nil, isVerified).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Add(user, isVerified)
	assert.NoError(t, err)
}

func TestUserRepository_Check(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := userRepo.NewUserRepository(mockDB)
	user := userEntity.UserLogin{
		Email:    "john@example.com",
		Password: "password123",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	mock.ExpectQuery("SELECT id, email, password, is_admin, is_verified FROM users WHERE email=\\$1").
		WithArgs(user.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "is_admin", "is_verified"}).
			AddRow(1, user.Email, hashedPassword, false, true))

	result, err := repo.Check(user)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	mockDB, mock := setupMockDB(t)
	defer mockDB.Close()

	repo := userRepo.NewUserRepository(mockDB)
	id := int32(1)

	mock.ExpectExec("DELETE FROM users WHERE id = \\$1").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(id)
	assert.NoError(t, err)
}
