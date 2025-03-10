package svc_test

import (
	"testing"

	questionEntity "github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/svc"
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

func (m *MockQuestionRepo) Add(q questionEntity.SetQuestion) error {
	args := m.Called(q)
	return args.Error(0)
}

func (m *MockQuestionRepo) Exists(setID int32, number int) (bool, error) {
	args := m.Called(setID, number)
	return args.Bool(0), args.Error(1)
}

func (m *MockQuestionRepo) List(filter map[string]string) ([]questionEntity.ListQuestionExample, error) {
	args := m.Called(filter)
	return args.Get(0).([]questionEntity.ListQuestionExample), args.Error(1)
}

func (m *MockQuestionRepo) Detail(id int32) (questionEntity.DetailQuestionExample, error) {
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

func (m *MockQuestionRepo) ListQuizQuestions(filter map[string]string) ([]questionEntity.ListQuestionQuiz, error) {
	args := m.Called(filter)
	return args.Get(0).([]questionEntity.ListQuestionQuiz), args.Error(1)
}

func (m *MockQuestionRepo) ListAdmin(filter map[string]string, page, limit int) ([]questionEntity.ListQuestionAdmin, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]questionEntity.ListQuestionAdmin), args.Error(1)
}

func (m *MockQuestionRepo) Edit(id int32, question questionEntity.EditQuestion) error {
	args := m.Called(id, question)
	return args.Error(0)
}

func (m *MockQuestionRepo) EditAnswer(id int32, answer questionEntity.EditAnswer) error {
	args := m.Called(id, answer)
	return args.Error(0)
}

func TestListAdminService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	mockData := []questionEntity.ListQuestionAdmin{
		{ID: 1, Number: 1, Type: "c4_faktual", Format: "mm", Content: "Question 1", IsQuiz: true, SetID: 1, SetName: "Set 1", LessonName: "Lesson 1", ClassName: "Class 1", },
	}

	mockRepo.On("ListAdmin", mock.Anything, 1, 10).Return(mockData, nil)

	questions, err := service.ListAdmin(map[string]string{}, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, questions, 1)
	assert.Equal(t, "Question 1", questions[0].Content)
}

func TestDetailQuestionService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	mockData := questionEntity.DetailQuestionExample{
		ID: 1, Number: 1, Type: "c4_faktual", Format: "mm", Content: "Question 1aaa", SetID: 9, Explanation: "exp-1",
	}

	mockRepo.On("Detail", mock.Anything).Return(mockData, nil)

	questions, err := service.DetailQuestion(1)
	assert.NoError(t, err)
	assert.Equal(t, "Question 1aaa", questions.Content)
}

func TestAddQuestionService(t *testing.T) {
	mockRepo := new(MockQuestionRepo)
	service := svc.NewQuestionService(mockRepo)

	question := questionEntity.SetQuestion{SetID: 1, Number: 1, Content: "New Question"}

	mockRepo.On("Exists", question.SetID, question.Number).Return(false, nil)
	mockRepo.On("Add", question).Return(nil)

	err := service.AddQuestion(question)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Exists", question.SetID, question.Number)
	mockRepo.AssertCalled(t, "Add", question)
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
		QuestionID: 8,
		Code:       "a",
		Content:    "Updated Question",
		ImgURL:     func(s string) *string { return &s }("http://random"),
		IsAnswer:   true,
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
