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

type AuthorHandler struct {
	authorService svc.AuthorService
	val           *validator.Validate
}

func NewAuthorHandler(authorService svc.AuthorService, val *validator.Validate) *AuthorHandler {
	return &AuthorHandler{authorService, val}
}

func (h *AuthorHandler) Router(r fiber.Router) {
	r.Post("/authors", m.R100(), h.AddAuthorHandler)
	r.Get("/authors", m.R100(), h.ListAuthorHandler)
	r.Get("/authors/:id", m.R100(), h.DetailAuthorHandler)
	// r.Put("/authors/:id", m.R100(), h.EditAuthorHandler)
	// r.Delete("/authors/:id", m.R100(), h.DeleteAuthorHandler)
}

func (h *AuthorHandler) AddAuthorHandler(c *fiber.Ctx) error {
	var author entity.CreateAuthorRequest
	if err := c.BodyParser(&author); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(author); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	if err := h.authorService.CreateAuthor(c.Context(), author); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "author added successfully", nil)
}

func (h *AuthorHandler) ListAuthorHandler(c *fiber.Ctx) error {
	data, err := h.authorService.ListAuthors(c.Context())
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendSuccess(c, "authors retrieved successfully", data)
}

func (h *AuthorHandler) DetailAuthorHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	author, err := h.authorService.GetAuthor(c.Context(), int32(id))
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "author detail retrieved successfully", author)
}

func (h *AuthorHandler) EditAuthorHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	var author entity.UpdateAuthorRequest
	if err := c.BodyParser(&author); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(author); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	if err := h.authorService.UpdateAuthor(c.Context(), id, author); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "author updated successfully", nil)
}

func (h *AuthorHandler) DeleteAuthorHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	if err := h.authorService.DeleteAuthor(c.Context(), int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "author deleted successfully", nil)
}
