package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ghulammuzz/misterblast/internal/question/entity"
	"github.com/ghulammuzz/misterblast/internal/question/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
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
	r.Post("/question", h.AddQuestionHandler)
	r.Put("/question/:id", h.EditQuestionHandler)
	r.Get("/question/:id", h.DetailQuestionsHandler)
	r.Get("/question", h.ListQuestionsHandler)
	r.Delete("/question/:id", h.DeleteQuestionHandler)

	// answer
	r.Delete("/answer/:id", h.DeleteAnswerHandler)
	r.Put("/answer/:id", h.EditAnswerHandler)
	r.Post("/quiz-answer", h.AddQuizAnswerHandler)

	// quiz
	r.Get("/quiz", h.ListQuizHandler)

	// admin
	r.Get("/admin-question", h.ListQuestionAdminHandler)
}

func (h *QuestionHandler) AddQuestionHandler(c *fiber.Ctx) error {
	var question entity.SetQuestion

	if err := c.BodyParser(&question); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(question); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}
	if err := h.questionService.AddQuestion(question); err != nil {
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

	questions, err := h.questionService.ListQuestions(filter)
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

	question, err := h.questionService.DetailQuestion(int32(id))
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
	var answer entity.SetAnswer

	if err := c.BodyParser(&answer); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(answer); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}
	if err := h.questionService.AddQuizAnswer(answer); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "answer added successfully", nil)
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

	questions, err := h.questionService.ListQuizQuestions(filter)
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

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	questions, err := h.questionService.ListAdmin(filter, page, limit)
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
