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

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/handler"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
)

type MockQuestionService struct {
	mock.Mock
}

func (m *MockQuestionService) AddQuestion(question questionEntity.SetQuestion) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockQuestionService) EditQuestion(id int32, question questionEntity.EditQuestion) error {
	args := m.Called(id, question)
	return args.Error(0)
}

func (m *MockQuestionService) ListQuestions(filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	args := m.Called(filter)
	return args.Get(0).([]questionEntity.ListQuestionExample), args.Error(1)
}

func (m *MockQuestionService) DeleteQuestion(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuestionService) DeleteAnswer(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuestionService) AddQuizAnswer(answer questionEntity.SetAnswer) error {
	args := m.Called(answer)
	return args.Error(0)
}

func (m *MockQuestionService) ListQuizQuestions(filter map[string]string) ([]questionEntity.ListQuestionQuiz, error) {
	args := m.Called(filter)
	return args.Get(0).([]questionEntity.ListQuestionQuiz), args.Error(1)
}

func (m *MockQuestionService) ListAdmin(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).(*response.PaginateResponse), args.Error(1)
}

func (m *MockQuestionService) DetailQuestion(id int32) (questionEntity.DetailQuestionExample, error) {
	args := m.Called(id)
	return args.Get(0).(questionEntity.DetailQuestionExample), args.Error(1)
}

func (m *MockQuestionService) EditQuizAnswer(id int32, answer questionEntity.EditAnswer) error {
	args := m.Called(id, answer)
	return args.Error(0)
}

func TestAddQuestionHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Post("/question", handler.AddQuestionHandler)

	question := questionEntity.SetQuestion{SetID: 9, Number: 1, Type: "c4_faktual", Format: "mm", Content: "Sample Question", Explanation: "exp-1"}
	questionJSON, _ := json.Marshal(question)

	mockService.On("AddQuestion", question).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/question", bytes.NewReader(questionJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestEditQuestionHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Put("/question/:id", handler.EditQuestionHandler)

	editQuestion := questionEntity.EditQuestion{SetID: 9, Number: 2, Type: "c3_faktual", Format: "mm", Content: "Updated Content", IsQuiz: false, Explanation: "exp-1"}
	editJSON, _ := json.Marshal(editQuestion)

	mockService.On("EditQuestion", int32(1), editQuestion).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/question/1", bytes.NewReader(editJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListQuestionsHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Get("/question", handler.ListQuestionsHandler)

	mockService.On("ListQuestions", mock.Anything).Return([]questionEntity.ListQuestionExample{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/question", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDetailQuestionsHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Get("/question/:id", handler.DetailQuestionsHandler)

	mockService.On("DetailQuestion", mock.Anything).Return(questionEntity.DetailQuestionExample{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/question/9", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteQuestionHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Delete("/question/:id", handler.DeleteQuestionHandler)

	mockService.On("DeleteQuestion", int32(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/question/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteAnswerHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Delete("/answer/:id", handler.DeleteAnswerHandler)

	mockService.On("DeleteAnswer", int32(11)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/answer/11", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestEditAnswerHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockQuestionService)
	validate := validator.New()
	handler := handler.NewQuestionHandler(mockService, validate)
	app.Put("/answer/:id", handler.EditAnswerHandler)

	editAnswer := questionEntity.EditAnswer{QuestionID: 8,
		Code:     "a",
		Content:  "Updated Question",
		ImgURL:   func(s string) *string { return &s }("http://random"),
		IsAnswer: true}
	editJSON, _ := json.Marshal(editAnswer)

	mockService.On("EditQuizAnswer", int32(1), editAnswer).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/answer/1", bytes.NewReader(editJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}
