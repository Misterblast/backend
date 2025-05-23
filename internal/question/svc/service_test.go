package svc_test

import (
	"context"
	"testing"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/svc"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo untuk menggantikan repo dalam pengujian
type MockQuestionRepo struct {
	mock.Mock
}

func (m *MockQuestionRepo) AddQuizAnswer(question questionEntity.SetAnswer) error {
	args := m.Called(question)
	return args.Error(0)
}

func (m *MockQuestionRepo) Add(q questionEntity.SetQuestion, lang string) error {
	args := m.Called(q, lang)
	return args.Error(0)
}

func (m *MockQuestionRepo) Exists(setID int32, number int) (bool, error) {
	args := m.Called(setID, number)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuestionRepo) List(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]questionEntity.ListQuestionExample), args.Error(1)
}

func (m *MockQuestionRepo) Detail(ctx context.Context, id int32) (questionEntity.DetailQuestionExample, error) {
	args := m.Called()
	return args.Get(0).(questionEntity.DetailQuestionExample), args.Error(1)
}

func (m *MockQuestionRepo) Delete(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuestionRepo) DeleteAnswer(id int32) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockQuestionRepo) ListQuizQuestions(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionQuiz, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]questionEntity.ListQuestionQuiz), args.Error(1)
}

func (m *MockQuestionRepo) ListQuizQuestionsLessonClass(ctx context.Context, filter map[string]string) ([]questionEntity.ListQuestionQuiz, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]questionEntity.ListQuestionQuiz), args.Error(1)
}

func (m *MockQuestionRepo) ListAdmin(ctx context.Context, filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
	args := m.Called(ctx, filter, page, limit)
	return args.Get(0).(*response.PaginateResponse), args.Error(1)
}

func (m *MockQuestionRepo) Edit(id int32, question questionEntity.EditQuestion) error {
	args := m.Called(id, question)
	return args.Error(0)
}

func (m *MockQuestionRepo) EditAnswer(id int32, answer questionEntity.EditAnswer) error {
	args := m.Called(id, answer)
	return args.Error(0)
}

func (m *MockQuestionRepo) ListQuestionTypes(ctx context.Context) ([]questionEntity.QuestionType, error) {
	args := m.Called(ctx)
	return args.Get(0).([]questionEntity.QuestionType), args.Error(1)
}

func (m *MockQuestionRepo) AddQuizAnswerBulk(answers []questionEntity.SetAnswer) error {
	args := m.Called(answers)
	return args.Error(0)
}

func TestListAdminService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	mockData := []questionEntity.ListQuestionAdmin{
		{ID: 1, Number: 1, Type: "c4_faktual", Format: "mm", Content: "Question 1", IsQuiz: true, SetID: 1, SetName: "Set 1", LessonName: "Lesson 1", ClassName: "Class 1"},
	}

	mockResponse := &response.PaginateResponse{
		Total: 1,
		Page:  1,
		Limit: 10,
		Data:  mockData,
	}

	mockRepo.On("ListAdmin", mock.Anything, mock.Anything, 1, 10).Return(mockResponse, nil)

	result, err := service.ListAdmin(context.Background(), map[string]string{}, 1, 10)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Total) // Mengecek total data
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.Limit)
	assert.Len(t, result.Data, 1)

	questions, ok := result.Data.([]questionEntity.ListQuestionAdmin)
	assert.True(t, ok)
	assert.Equal(t, "Question 1", questions[0].Content)

	mockRepo.AssertExpectations(t)
}

func TestDetailQuestionService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	mockData := questionEntity.DetailQuestionExample{
		ID: 1, Number: 1, Type: "c4_faktual", Format: "mm", Content: "Question 1aaa", SetID: 9, Explanation: "exp-1",
	}

	mockRepo.On("Detail", mock.Anything).Return(mockData, nil)

	questions, err := service.DetailQuestion(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Question 1aaa", questions.Content)
}

func TestAddQuestionService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	question := questionEntity.SetQuestion{
		SetID:   1,
		Number:  1,
		Content: "New Question",
	}

	mockRepo.On("Exists", question.SetID, question.Number).Return(false, nil)
	mockRepo.On("Add", question, "en").Return(nil)

	err := service.AddQuestion(question, "en")
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "Exists", question.SetID, question.Number)
	mockRepo.AssertCalled(t, "Add", question, "en")
}

func TestDeleteQuestionService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	mockRepo.On("Delete", int32(1)).Return(nil)

	err := service.DeleteQuestion(1)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", int32(1))
}

func TestEditQuestionService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	question := questionEntity.EditQuestion{
		Number:      1,
		Type:        "c4_faktual",
		Format:      "mm",
		Content:     "Updated Question",
		IsQuiz:      false,
		SetID:       1,
		Explanation: "exp-1",
	}

	mockRepo.On("Edit", int32(1), question).Return(nil)

	err := service.EditQuestion(1, question)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Edit", int32(1), question)
}

func TestEditAnswerService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	answer := questionEntity.EditAnswer{
		Code:     "a",
		Content:  "Updated Question",
		ImgURL:   func(s string) *string { return &s }("http://random"),
		IsAnswer: true,
	}

	mockRepo.On("EditAnswer", int32(1), answer).Return(nil)

	err := service.EditQuizAnswer(1, answer)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "EditAnswer", int32(1), answer)
}

func TestDeleteAnswerService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	mockRepo.On("DeleteAnswer", int32(8)).Return(nil)

	err := service.DeleteAnswer(8)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "DeleteAnswer", int32(8))
}
