package handler

import (
	"strconv"

	"github.com/ghulammuzz/misterblast/internal/content/entity"
	"github.com/ghulammuzz/misterblast/internal/content/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ContentHandler struct {
	contentService svc.ContentService
	val            *validator.Validate
}

func NewContentHandler(contentService svc.ContentService, val *validator.Validate) *ContentHandler {
	return &ContentHandler{contentService, val}
}

func (h *ContentHandler) Router(r fiber.Router) {
	r.Post("/content", m.R100(), h.AddContentHandler)
	r.Get("/content", m.R100(), h.ListContentHandler)
	r.Get("/content/:id", m.R100(), h.DetailContentHandler)
	r.Put("/content/:id", m.R100(), h.EditContentHandler)
	r.Delete("/content/:id", m.R100(), h.DeleteContentHandler)
}

func (h *ContentHandler) AddContentHandler(c *fiber.Ctx) error {
	var content entity.Content
	if err := c.BodyParser(&content); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(content); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	lang := c.Query("lang")
	if lang != "id" && lang != "en" {
		return response.SendError(c, fiber.StatusBadRequest, "invalid lang parameter", nil)
	}

	if err := h.contentService.Add(content, lang); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "content added successfully", nil)
}

func (h *ContentHandler) ListContentHandler(c *fiber.Ctx) error {
	filter := map[string]string{
		"lang":  c.Query("lang"),
		"page":  c.Query("page"),
		"limit": c.Query("limit"),
	}

	data, err := h.contentService.List(c.Context(), filter)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "content list retrieved successfully", data)
}

func (h *ContentHandler) DetailContentHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	data, err := h.contentService.Detail(c.Context(), int32(id))
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendSuccess(c, "content detail retrieved successfully", data)
}

func (h *ContentHandler) EditContentHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	var content entity.Content
	if err := c.BodyParser(&content); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(content); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.contentService.Edit(int32(id), content); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "content updated successfully", nil)
}

func (h *ContentHandler) DeleteContentHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	if err := h.contentService.Delete(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "content deleted successfully", nil)
}
