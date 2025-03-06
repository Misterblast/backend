package handler

import (
	"github.com/ghulammuzz/misterblast/internal/set/entity"
	"github.com/ghulammuzz/misterblast/internal/set/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type SetHandler struct {
	setService svc.SetService
	val        *validator.Validate
}

func NewSetHandler(setService svc.SetService, val *validator.Validate) *SetHandler {
	return &SetHandler{setService, val}
}

func (h *SetHandler) Router(r fiber.Router) {
	r.Post("/set", h.AddSetHandler)
	r.Delete("/set/:id", h.DeleteSetHandler)
	r.Get("/set", h.ListSetsHandler)
}

func (h *SetHandler) AddSetHandler(c *fiber.Ctx) error {
	var set entity.SetSet

	if err := c.BodyParser(&set); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(set); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.setService.AddSet(set); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "set added successfully", nil)
}

func (h *SetHandler) DeleteSetHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	if err := h.setService.DeleteSet(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "set deleted successfully", nil)
}

func (h *SetHandler) ListSetsHandler(c *fiber.Ctx) error {
	filter := map[string]string{}
	if class := c.Query("class"); class != "" {
		filter["class"] = class
	}
	if lesson := c.Query("lesson"); lesson != "" {
		filter["lesson"] = lesson
	}
	if isQuiz := c.Query("is_quiz"); isQuiz != "" {
		filter["is_quiz"] = isQuiz
	}

	sets, err := h.setService.ListSets(filter)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "sets retrieved successfully", sets)
}
