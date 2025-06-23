package handler

import (
	"github.com/ghulammuzz/misterblast/internal/quiz/entity"
	"github.com/ghulammuzz/misterblast/internal/quiz/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type QuizHandler struct {
	quizService svc.QuizService
	val         *validator.Validate
}

func NewQuizHandler(quizService svc.QuizService, val *validator.Validate) *QuizHandler {
	return &QuizHandler{quizService, val}
}

func (h *QuizHandler) Router(r fiber.Router) {
	r.Post("/submit-quiz/:set_id", m.JWTProtected(), m.R100(), h.SubmitQuizHandler)

	r.Get("/quiz-submission-admin", m.R100(), h.AdminQuizSubmissionHandler)
	r.Get("/quiz-submission", m.JWTProtected(), m.R100(), h.QuizSubmissionHandler)
	r.Get("/quiz-submission/:submission_id", m.JWTProtected(), m.R100(), h.GetSubmissionDetailHandler)
	r.Get("/quiz-result", m.JWTProtected(), m.R100(), h.GetResultHandler)

}

func (h *QuizHandler) SubmitQuizHandler(c *fiber.Ctx) error {
	var req entity.QuizSubmit

	setID, err := c.ParamsInt("set_id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	lang := c.Get("Lang")
	if lang == "" {
		lang = c.Query("lang")
	}
	if lang == "" {
		lang = "id"
	}

	if err := c.BodyParser(&req); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(req); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", err.Error())
	}
	id, err := h.quizService.SubmitQuiz(req, setID, userID, lang)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "question added successfully", id)
}

func (h *QuizHandler) AdminQuizSubmissionHandler(c *fiber.Ctx) error {
	filter := map[string]string{}

	if c.Query("class_id") != "" {
		filter["class_id"] = c.Query("class_id")
	}
	if c.Query("lesson_id") != "" {
		filter["lesson_id"] = c.Query("lesson_id")
	}

	if c.Query("type") != "" {
		filter["type"] = c.Query("type")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	quiz, err := h.quizService.ListAdmin(filter, page, limit)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "quiz admin retrieved successfully", quiz)
}

func (h *QuizHandler) QuizSubmissionHandler(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	filter := map[string]string{}

	if c.Query("class_id") != "" {
		filter["class_id"] = c.Query("class_id")
	}
	if c.Query("lesson_id") != "" {
		filter["lesson_id"] = c.Query("lesson_id")
	}
	if c.Query("type") != "" {
		filter["type"] = c.Query("type")
	}

	if c.Query("page") != "" {
		filter["page"] = c.Query("page")
	} else {
		filter["page"] = "1"
	}
	if c.Query("limit") != "" {
		filter["limit"] = c.Query("limit")
	} else {
		filter["limit"] = "10"
	}
	quiz, err := h.quizService.List(filter, userID)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "quiz admin retrieved successfully", quiz)
}

func (h *QuizHandler) GetResultHandler(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	quiz, err := h.quizService.GetResult(userID)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "quiz result retrieved successfully", quiz)
}

func (h *QuizHandler) GetSubmissionDetailHandler(c *fiber.Ctx) error {
	submissionId, err := c.ParamsInt("submission_id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid submission ID", nil)
	}
	submission, err := h.quizService.GetSubmissionResult(submissionId)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "quiz submission detail retrieved successfully", submission)
}
