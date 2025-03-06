package svc_test

import (
	"errors"
	"testing"

	classEntity "github.com/ghulammuzz/misterblast/internal/class/entity"
	"github.com/ghulammuzz/misterblast/internal/class/svc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Exists(name string) (bool, error) {
	args := m.Called(name)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) Add(class classEntity.SetClass) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockRepo) Delete(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) List() ([]classEntity.Class, error) {
	args := m.Called()
	return args.Get(0).([]classEntity.Class), args.Error(1)
}

func TestAddClass(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := svc.NewClassService(mockRepo)

	t.Run("should return error when name is empty", func(t *testing.T) {
		class := classEntity.SetClass{Name: ""}
		err := svc.AddClass(class)
		assert.EqualError(t, err, "name is required")
	})

	t.Run("should return error when Exists fails", func(t *testing.T) {
		mockRepo.On("Exists", "Math").Return(false, errors.New("db error")).Once()
		class := classEntity.SetClass{Name: "Math"}
		err := svc.AddClass(class)
		assert.EqualError(t, err, "db error")
	})

	t.Run("should return error when class already exists", func(t *testing.T) {
		mockRepo.On("Exists", "Math").Return(true, nil).Once()
		class := classEntity.SetClass{Name: "Math"}
		err := svc.AddClass(class)
		assert.EqualError(t, err, "lesson already exists")
	})

	t.Run("should return error when Add fails", func(t *testing.T) {
		mockRepo.On("Exists", "Math").Return(false, nil).Once()
		mockRepo.On("Add", mock.Anything).Return(errors.New("insert error")).Once()
		class := classEntity.SetClass{Name: "Math"}
		err := svc.AddClass(class)
		assert.EqualError(t, err, "insert error")
	})

	t.Run("should add class successfully", func(t *testing.T) {
		mockRepo.On("Exists", "Math").Return(false, nil).Once()
		mockRepo.On("Add", mock.Anything).Return(nil).Once()
		class := classEntity.SetClass{Name: "Math"}
		err := svc.AddClass(class)
		assert.NoError(t, err)
	})
}

func TestDeleteClass(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := svc.NewClassService(mockRepo)

	t.Run("should return error when id is invalid", func(t *testing.T) {
		err := svc.DeleteClass(0)
		assert.EqualError(t, err, "invalid id")
	})

	t.Run("should return error when Delete fails", func(t *testing.T) {
		mockRepo.On("Delete", int32(1)).Return(errors.New("delete error")).Once()
		err := svc.DeleteClass(1)
		assert.EqualError(t, err, "delete error")
	})

	t.Run("should delete class successfully", func(t *testing.T) {
		mockRepo.On("Delete", int32(1)).Return(nil).Once()
		err := svc.DeleteClass(1)
		assert.NoError(t, err)
	})
}

func TestListClasses(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := svc.NewClassService(mockRepo)

	t.Run("should return error when List fails", func(t *testing.T) {
		mockRepo.On("List").Return(([]classEntity.Class)(nil), errors.New("list error")).Once()
		_, err := svc.ListClasses()
		assert.EqualError(t, err, "list error")
	})

	t.Run("should list classes successfully", func(t *testing.T) {
		mockRepo.On("List").Return([]classEntity.Class{{ID: 1, Name: "Math"}}, nil).Once()
		classes, err := svc.ListClasses()
		assert.NoError(t, err)
		assert.Len(t, classes, 1)
		assert.Equal(t, "Math", classes[0].Name)
	})
}
