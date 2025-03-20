package handler

import (
	"errors"
	"net/http"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	service "github.com/ghulammuzz/misterblast/internal/task/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TaskHandler struct {
	s   service.TaskService
	val *validator.Validate
}

func NewTaskHandler(s service.TaskService, val *validator.Validate) *TaskHandler {
	return &TaskHandler{s, val}
}

func (h *TaskHandler) Router(r fiber.Router) {
	r.Get("/tasks", m.R100(), m.JWTProtected(), h.List)
	r.Post("/tasks", m.R100(), m.JWTProtected(), h.CreateTask)
}

func (h *TaskHandler) List(c *fiber.Ctx) error {
	var listTaskRequestDto entity.ListTaskRequestDto
	listTaskRequestDto.Page = int32(c.QueryInt("page", 1))
	listTaskRequestDto.Search = c.Query("search")

	tasks, err := h.s.List(listTaskRequestDto)
	if err != nil {
		var appErr *app.AppError
		if errors.As(err, &appErr) {
			return response.SendError(c, appErr.Code, appErr.Message, appErr.Error())
		}
		return response.SendError(c, fiber.StatusInternalServerError, "failed to get tasks", nil)
	}

	return response.SendSuccess(c, "tasks retrieved successfully", tasks)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var createTaskRequestDto entity.CreateTaskRequestDto

	if err := c.BodyParser(&createTaskRequestDto); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}
	form, err := c.MultipartForm()
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid multipart form", nil)
	}
	if files, ok := form.File["attachments"]; ok {
		createTaskRequestDto.Attachments = files
	}

	if err := h.val.Struct(createTaskRequestDto); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	claims := c.Locals("claims").(jwt.MapClaims)

	if err := h.s.Create(int32(claims["user_id"].(float64)), createTaskRequestDto); err != nil {
		var appErr *app.AppError
		ok := errors.As(err, &appErr)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendResponse(c, http.StatusCreated, "Task added successfully", nil)
}
