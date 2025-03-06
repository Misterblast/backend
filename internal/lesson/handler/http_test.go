package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ghulammuzz/misterblast/internal/lesson/entity"
	"github.com/ghulammuzz/misterblast/internal/lesson/handler"
)

type MockLessonService struct {
	mock.Mock
}

func (m *MockLessonService) AddLesson(lesson entity.Lesson) error {
	args := m.Called(lesson)
	return args.Error(0)
}

func (m *MockLessonService) DeleteLesson(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockLessonService) ListLessons() ([]entity.Lesson, error) {
	args := m.Called()
	return args.Get(0).([]entity.Lesson), args.Error(1)
}

func TestAddLessonHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockLessonService)
	validator := validator.New()
	h := handler.NewLessonHandler(mockService, validator)
	app.Post("/lesson", h.AddLessonHandler)

	validLesson := entity.Lesson{Name: "Sample Lesson"}
	mockService.On("AddLesson", validLesson).Return(nil)

	body, _ := json.Marshal(validLesson)
	req := httptest.NewRequest(http.MethodPost, "/lesson", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteLessonHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockLessonService)
	validator := validator.New()
	h := handler.NewLessonHandler(mockService, validator)
	app.Delete("/lesson/:id", h.DeleteLessonHandler)

	mockService.On("DeleteLesson", int32(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/lesson/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListLessonsHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockLessonService)
	validator := validator.New()
	h := handler.NewLessonHandler(mockService, validator)
	app.Get("/lesson", h.ListLessonsHandler)

	mockLessons := []entity.Lesson{
		{ID: 1, Name: "Lesson 1"},
		{ID: 2, Name: "Lesson 2"},
	}
	mockService.On("ListLessons").Return(mockLessons, nil)

	req := httptest.NewRequest(http.MethodGet, "/lesson", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
