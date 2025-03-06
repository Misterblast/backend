package handler

import (
	classEntity "github.com/ghulammuzz/misterblast/internal/class/entity"
	classSvc "github.com/ghulammuzz/misterblast/internal/class/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type ClassHandler struct {
	classService classSvc.ClassService
}

func NewClassHandler(classService classSvc.ClassService) *ClassHandler {
	return &ClassHandler{classService: classService}
}

func (h *ClassHandler) Router(r fiber.Router) {
	r.Post("/class", h.AddClassHandler)
	r.Delete("/class/:id", h.DeleteClassHandler)
	r.Get("/class", h.ListClassesHandler)
}

func (h *ClassHandler) AddClassHandler(c *fiber.Ctx) error {
	var newClass classEntity.SetClass

	if err := c.BodyParser(&newClass); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.classService.AddClass(newClass); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "class added successfully", nil)
}

func (h *ClassHandler) DeleteClassHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	if err := h.classService.DeleteClass(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "class deleted successfully", nil)
}

func (h *ClassHandler) ListClassesHandler(c *fiber.Ctx) error {
	classes, err := h.classService.ListClasses()
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "classes retrieved successfully", classes)
}
