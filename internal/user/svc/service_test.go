package svc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	userSvc "github.com/ghulammuzz/misterblast/internal/user/svc"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Exists(id int32) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) Detail(id int32) (userEntity.DetailUser, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return userEntity.DetailUser{}, args.Error(1)
	}
	return args.Get(0).(userEntity.DetailUser), args.Error(1)
}

func (m *MockUserRepository) AdminActivation(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Add(user userEntity.Register, isVerified bool) error {
	args := m.Called(user, isVerified)
	return args.Error(0)
}

func (m *MockUserRepository) Check(user userEntity.UserLogin) (*userEntity.UserJWT, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.UserJWT), args.Error(1)
}

func (m *MockUserRepository) Delete(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Auth(id int32) (userEntity.UserAuth, error) {
	args := m.Called(id)
	return args.Get(0).(userEntity.UserAuth), args.Error(1)
}

func (m *MockUserRepository) Edit(id int32, user userEntity.EditUser) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetIDByEmail(email string) (int32, error) {
	args := m.Called(email)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockUserRepository) List(filter map[string]string, page, limit int) ([]userEntity.ListUser, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]userEntity.ListUser), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	user := userEntity.Register{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	mockRepo.On("Add", user, true).Return(nil)
	err := service.Register(user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	user := userEntity.UserLogin{
		Email:    "john@example.com",
		Password: "password123",
	}

	userJWT := &userEntity.UserJWT{
		ID:         1,
		Email:      user.Email,
		IsAdmin:    false,
		IsVerified: true,
	}

	mockRepo.On("Check", user).Return(userJWT, nil)

	resp, token, err := service.Login(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Email, resp.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	id := int32(1)
	mockRepo.On("Delete", id).Return(nil)

	err := service.DeleteUser(id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_AuthUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	id := int32(1)
	userAuth := userEntity.UserAuth{
		ID:         id,
		Name:       "John Doe",
		Email:      "john@example.com",
		ImgUrl:     "",
		IsAdmin:    false,
		IsVerified: true,
	}

	mockRepo.On("Auth", id).Return(userAuth, nil)

	resp, err := service.AuthUser(id)
	assert.NoError(t, err)
	assert.Equal(t, userAuth, resp)
	mockRepo.AssertExpectations(t)
}
func TestUserService_ListUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	filter := map[string]string{"role": "user"}
	page, limit := 1, 10
	mockUsers := []userEntity.ListUser{{ID: 1, Name: "John Doe", Email: "john@example.com"}}

	mockRepo.On("List", filter, page, limit).Return(mockUsers, nil)

	resp, err := service.ListUser(filter, page, limit)
	assert.NoError(t, err)
	assert.Equal(t, mockUsers, resp)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DetailUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	id := int32(1)
	mockUser := userEntity.DetailUser{ID: id, Name: "John Doe", Email: "john@example.com"}

	mockRepo.On("Detail", id).Return(mockUser, nil)

	resp, err := service.DetailUser(id)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, resp)
	mockRepo.AssertExpectations(t)
}

func TestUserService_EditUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := userSvc.NewUserService(mockRepo)

	id := int32(1)
	userEdit := userEntity.EditUser{Name: "John Updated"}

	mockRepo.On("Edit", id, userEdit).Return(nil)

	err := service.EditUser(id, userEdit)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
