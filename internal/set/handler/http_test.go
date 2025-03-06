package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ghulammuzz/misterblast/internal/set/entity"
	"github.com/ghulammuzz/misterblast/internal/set/handler"
)

// Mock Service
type MockSetService struct {
	mock.Mock
}

func (m *MockSetService) AddSet(set entity.SetSet) error {
	args := m.Called(set)
	return args.Error(0)
}

func (m *MockSetService) DeleteSet(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSetService) ListSets(filter map[string]string) ([]entity.ListSet, error) {
	args := m.Called(filter)
	return args.Get(0).([]entity.ListSet), args.Error(1)
}

func TestAddSetHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockSetService)
	validate := validator.New()
	h := handler.NewSetHandler(mockService, validate)

	app.Post("/set", h.AddSetHandler)

	set := entity.SetSet{Name: "Set A", LessonID: 1, ClassID: 1}
	mockService.On("AddSet", set).Return(nil)

	body, _ := json.Marshal(set)
	req := httptest.NewRequest("POST", "/set", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteSetHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockSetService)
	validate := validator.New()
	h := handler.NewSetHandler(mockService, validate)

	app.Delete("/set/:id", h.DeleteSetHandler)

	mockService.On("DeleteSet", int32(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/set/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListSetsHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockSetService)
	validate := validator.New()
	h := handler.NewSetHandler(mockService, validate)

	app.Get("/set", h.ListSetsHandler)

	mockSets := []entity.ListSet{
		{ID: 1, Name: "Set A", Lesson: "Math", Class: "Class 1"},
		{ID: 2, Name: "Set B", Lesson: "Science", Class: "Class 2"},
	}
	mockService.On("ListSets", mock.Anything).Return(mockSets, nil)

	req := httptest.NewRequest("GET", "/set", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListSetsHandler_Error(t *testing.T) {
	app := fiber.New()
	mockService := new(MockSetService)
	validate := validator.New()
	h := handler.NewSetHandler(mockService, validate)

	app.Get("/set", h.ListSetsHandler)

	mockService.On("ListSets", mock.Anything).Return([]entity.ListSet{}, errors.New("database error"))

	req := httptest.NewRequest("GET", "/set", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
	mockService.AssertExpectations(t)
}
