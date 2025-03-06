package handler

import (
	"github.com/ghulammuzz/misterblast/internal/lesson/entity"
	"github.com/ghulammuzz/misterblast/internal/lesson/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LessonHandler struct {
	lessonService svc.LessonService
	val           *validator.Validate
}

func NewLessonHandler(lessonService svc.LessonService, val *validator.Validate) *LessonHandler {
	return &LessonHandler{lessonService, val}
}

func (h *LessonHandler) Router(r fiber.Router) {
	r.Post("/lesson", h.AddLessonHandler)
	r.Delete("/lesson/:id", h.DeleteLessonHandler)
	r.Get("/lesson", h.ListLessonsHandler)
}

func (h *LessonHandler) AddLessonHandler(c *fiber.Ctx) error {
	var lesson entity.Lesson

	if err := c.BodyParser(&lesson); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	if err := h.val.Struct(lesson); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.lessonService.AddLesson(lesson); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "lesson added successfully", nil)
}

func (h *LessonHandler) DeleteLessonHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	if err := h.lessonService.DeleteLesson(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "lesson deleted successfully", nil)
}

func (h *LessonHandler) ListLessonsHandler(c *fiber.Ctx) error {
	lessons, err := h.lessonService.ListLessons()
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "lessons retrieved successfully", lessons)
}
