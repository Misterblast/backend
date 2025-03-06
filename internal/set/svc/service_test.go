package svc_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ghulammuzz/misterblast/internal/set/entity"
	"github.com/ghulammuzz/misterblast/internal/set/svc"
)

// Mock Repository
type MockSetRepository struct {
	mock.Mock
}

func (m *MockSetRepository) Add(set entity.SetSet) error {
	args := m.Called(set)
	return args.Error(0)
}

func (m *MockSetRepository) Delete(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSetRepository) List(filter map[string]string) ([]entity.ListSet, error) {
	args := m.Called(filter)
	return args.Get(0).([]entity.ListSet), args.Error(1)
}

func TestAddSet(t *testing.T) {
	mockRepo := new(MockSetRepository)
	service := svc.NewSetService(mockRepo)

	set := entity.SetSet{Name: "Set A", LessonID: 1, ClassID: 1}
	mockRepo.On("Add", set).Return(nil)

	err := service.AddSet(set)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteSet(t *testing.T) {
	mockRepo := new(MockSetRepository)
	service := svc.NewSetService(mockRepo)

	mockRepo.On("Delete", int32(1)).Return(nil)

	err := service.DeleteSet(1)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListSets(t *testing.T) {
	mockRepo := new(MockSetRepository)
	service := svc.NewSetService(mockRepo)

	mockSets := []entity.ListSet{
		{ID: 1, Name: "Set A", Lesson: "Math", Class: "Class 1"},
		{ID: 2, Name: "Set B", Lesson: "Science", Class: "Class 2"},
	}
	mockRepo.On("List", mock.Anything).Return(mockSets, nil)

	sets, err := service.ListSets(map[string]string{})
	assert.NoError(t, err)
	assert.Len(t, sets, 2)
	mockRepo.AssertExpectations(t)
}

func TestListSets_Error(t *testing.T) {
	mockRepo := new(MockSetRepository)
	service := svc.NewSetService(mockRepo)

	mockRepo.On("List", mock.Anything).Return([]entity.ListSet{}, errors.New("database error"))

	sets, err := service.ListSets(map[string]string{})
	assert.Error(t, err)
	assert.Empty(t, sets)
	mockRepo.AssertExpectations(t)
}
