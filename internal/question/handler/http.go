package handler

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type QuestionHandler struct {
	questionService svc.QuestionService
	val             *validator.Validate
}

func NewQuestionHandler(questionService svc.QuestionService, val *validator.Validate) *QuestionHandler {
	return &QuestionHandler{questionService, val}
}

func (h *QuestionHandler) Router(r fiber.Router) {
	// question
	r.Post("/question", m.R100(), h.AddQuestionHandler)
	r.Put("/question/:id", m.R100(), h.EditQuestionHandler)
	r.Get("/question/:id", m.R100(), h.DetailQuestionsHandler)
	r.Get("/question", m.R100(), h.ListQuestionsHandler)
	r.Delete("/question/:id", m.R100(), h.DeleteQuestionHandler)

	// answer
	r.Delete("/answer/:id", m.R100(), h.DeleteAnswerHandler)
	r.Put("/answer/:id", m.R100(), h.EditAnswerHandler)
	r.Post("/quiz-answer", m.R100(), h.AddQuizAnswerHandler)
	r.Post("/question-answer", m.R100(), h.AddQuizAnswerHandler)
	r.Post("/quiz-answer-bulk/:id", m.R100(), h.AddQuizAnswerBulkHandler)
	r.Post("/question-answer-bulk/:id", m.R100(), h.AddQuizAnswerBulkHandler)

	// quiz
	r.Get("/quiz", m.R100(), h.ListQuizHandler)

	// admin
	r.Get("/admin-question", m.R100(), h.ListQuestionAdminHandler)

	// question type
	r.Get("/question-type", m.R100(), h.ListQuestionTypes)
}

func (h *QuestionHandler) AddQuestionHandler(c *fiber.Ctx) error {
	var question entity.SetQuestion

	if err := c.BodyParser(&question); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(question); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	lang := c.Query("lang")
	if lang == "" {
		return response.SendError(c, fiber.StatusBadRequest, "language (lang) is required", nil)
	}

	if err := h.questionService.AddQuestion(question, lang); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "question added successfully", nil)
}

func (h *QuestionHandler) ListQuestionsHandler(c *fiber.Ctx) error {
	filter := map[string]string{}
	if setID := c.Query("set_id"); setID != "" {
		filter["set_id"] = setID
	}
	if lessonID := c.Query("lesson_id"); lessonID != "" {
		filter["lesson_id"] = lessonID
	}
	if classID := c.Query("class_id"); classID != "" {
		filter["class_id"] = classID
	}
	if isQuiz := c.Query("is_quiz"); isQuiz != "" {
		filter["is_quiz"] = isQuiz
	}

	questions, err := h.questionService.ListQuestions(c.Context(), filter)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "questions retrieved successfully", questions)
}

func (h *QuestionHandler) DetailQuestionsHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid question ID", nil)
	}

	question, err := h.questionService.DetailQuestion(c.Context(), int32(id))
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "question retrieved successfully", question)
}

func (h *QuestionHandler) DeleteQuestionHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid question ID", nil)
	}

	if err := h.questionService.DeleteQuestion(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "question deleted successfully", nil)
}

// Quiz Answer

func (h *QuestionHandler) AddQuizAnswerHandler(c *fiber.Ctx) error {
	var answers []entity.SetAnswer

	questionID, err := c.ParamsInt("id")
	if err != nil || questionID <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid question ID", nil)
	}
	if err := c.BodyParser(&answers); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	for i, ans := range answers {
		if err := h.val.Struct(ans); err != nil {
			return response.SendError(c, fiber.StatusBadRequest,
				fmt.Sprintf("Validation failed at index %d", i),
				err.Error())
		}
	}

	if err := h.questionService.AddQuizAnswerBulk(int32(questionID), answers); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendSuccess(c, "answers added successfully", nil)
}

func (h *QuestionHandler) AddQuizAnswerBulkHandler(c *fiber.Ctx) error {
	var answers []entity.SetAnswer

	questionID, err := c.ParamsInt("id")
	if err != nil || questionID <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid question ID", nil)
	}
	if err := c.BodyParser(&answers); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	for i, ans := range answers {
		if err := h.val.Struct(ans); err != nil {
			return response.SendError(c, fiber.StatusBadRequest,
				fmt.Sprintf("Validation failed at index %d", i),
				err.Error())
		}
	}

	if err := h.questionService.AddQuizAnswerBulk(int32(questionID), answers); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendSuccess(c, "answers added successfully", nil)
}

func (h *QuestionHandler) DeleteAnswerHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid answer ID", nil)
	}

	if err := h.questionService.DeleteAnswer(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "answer deleted successfully", nil)
}

// Quiz

func (h *QuestionHandler) ListQuizHandler(c *fiber.Ctx) error {
	filter := map[string]string{}
	if c.Query("set_id") != "" {
		filter["set_id"] = c.Query("set_id")
	}
	if c.Query("type") != "" {
		filter["type"] = c.Query("type")
	}
	if c.Query("number") != "" {
		filter["number"] = c.Query("number")
	}

	if c.Query("class_id") != "" {
		filter["class_id"] = c.Query("class_id")
	}
	if c.Query("lesson_id") != "" {
		filter["lesson_id"] = c.Query("lesson_id")
	}
	lang := c.Get("Lang")
	if lang == "" {
		lang = "id"
	} else if lang != "id" && lang != "en" {
		return response.SendError(c, 400, "invalid lang, only 'id' or 'en' allowed", nil)
	}
	filter["lang"] = lang

	questions, err := h.questionService.ListQuizQuestions(c.Context(), filter)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "questions retrieved successfully", questions)
}

// admin
func (h *QuestionHandler) ListQuestionAdminHandler(c *fiber.Ctx) error {
	filter := map[string]string{}
	if c.Query("is_quiz") != "" {
		filter["is_quiz"] = c.Query("is_quiz")
	}
	if c.Query("lesson") != "" {
		filter["lesson"] = c.Query("lesson")
	}
	if c.Query("class") != "" {
		filter["class"] = c.Query("class")
	}
	if c.Query("set") != "" {
		filter["set"] = c.Query("set")
	}
	if c.Query("search") != "" {
		filter["search"] = c.Query("search")
	}
	if c.Query("lessonCode") != "" {
		filter["code"] = c.Query("lessonCode")
	}
	lang := c.Get("Lang")
	if lang == "" {
		lang = "id"
	} else if lang != "id" && lang != "en" {
		return response.SendError(c, 400, "invalid lang, only 'id' or 'en' allowed", nil)
	}
	filter["lang"] = lang

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	questions, err := h.questionService.ListAdmin(c.Context(), filter, page, limit)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "questions admin retrieved successfully", questions)
}

func (h *QuestionHandler) EditQuestionHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid question ID", nil)
	}

	var question entity.EditQuestion
	if err := c.BodyParser(&question); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(question); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	if err := h.questionService.EditQuestion(int32(id), question); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "question updated successfully", nil)
}

func (h *QuestionHandler) EditAnswerHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return response.SendError(c, fiber.StatusBadRequest, "invalid answer ID", nil)
	}

	var answer entity.EditAnswer
	if err := c.BodyParser(&answer); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(answer); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}

	if err := h.questionService.EditQuizAnswer(int32(id), answer); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "answer updated successfully", nil)
}

// q type
func (h *QuestionHandler) ListQuestionTypes(c *fiber.Ctx) error {
	questionTypes, err := h.questionService.ListQuestionTypes(c.Context())
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "question types retrieved successfully", questionTypes)
}
