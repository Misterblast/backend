package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	classEntity "github.com/ghulammuzz/misterblast/internal/class/entity"
	"github.com/ghulammuzz/misterblast/internal/class/handler"
)

type MockClassService struct {
	mock.Mock
}

func (m *MockClassService) AddClass(class classEntity.SetClass) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockClassService) DeleteClass(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockClassService) ListClasses() ([]classEntity.Class, error) {
	args := m.Called()
	return args.Get(0).([]classEntity.Class), args.Error(1)
}

func TestAddClassHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockClassService)
	handler := handler.NewClassHandler(mockService)
	app.Post("/class", handler.AddClassHandler)

	validClass := classEntity.SetClass{Name: "1"}
	mockService.On("AddClass", validClass).Return(nil)

	body, _ := json.Marshal(validClass)
	req := httptest.NewRequest(http.MethodPost, "/class", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteClassHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockClassService)
	handler := handler.NewClassHandler(mockService)
	app.Delete("/class/:id", handler.DeleteClassHandler)

	mockService.On("DeleteClass", int32(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/class/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListClassesHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockClassService)
	handler := handler.NewClassHandler(mockService)
	app.Get("/class", handler.ListClassesHandler)

	mockClasses := []classEntity.Class{
		{ID: 1, Name: "1"},
		{ID: 2, Name: "2"},
	}
	mockService.On("ListClasses").Return(mockClasses, nil)

	req := httptest.NewRequest(http.MethodGet, "/class", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
