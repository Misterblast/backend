package handler

import (
	"errors"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	service "github.com/ghulammuzz/misterblast/internal/task/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	s   service.TaskService
	val *validator.Validate
}

func NewTaskHandler(s service.TaskService, val *validator.Validate) *TaskHandler {
	return &TaskHandler{s, val}
}

func (h *TaskHandler) Router(r fiber.Router) {
	r.Get("/tasks", m.R100(), h.List)
	r.Get("/tasks/:id", m.R100(), h.Index)
	r.Post("/tasks", m.R100(), h.CreateTask)
	r.Delete("/tasks/:id", m.R100(), h.Delete)
	r.Post("/submit-task/:id", m.R100(), m.JWTProtected(), h.SubmitTask)
}

func (h *TaskHandler) List(c *fiber.Ctx) error {
	filter := map[string]string{}
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if c.Query("search") != "" {
		filter["search"] = c.Query("search")
	}

	tasks, err := h.s.List(filter, page, limit)
	if err != nil {
		var appErr *app.AppError
		ok := errors.As(err, &appErr)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "tasks retrieved successfully", tasks)
}

func (h *TaskHandler) Index(c *fiber.Ctx) error {
	taskId, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Params", nil)
	}
	task, err := h.s.Index(int32(taskId))
	if err != nil {
		var appErr *app.AppError
		ok := errors.As(err, &appErr)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendResponse(c, fiber.StatusOK, "Task retrieved", task)
}
func (h *TaskHandler) Delete(c *fiber.Ctx) error {
	taskId, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Params", nil)
	}
	err = h.s.Delete(int32(taskId))
	if err != nil {
		var appErr *app.AppError
		ok := errors.As(err, &appErr)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendResponse(c, fiber.StatusOK, "Task Deleted", nil)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	var createTaskRequestDto entity.CreateTaskRequestDto

	if err := c.BodyParser(&createTaskRequestDto); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid request body", nil)
	}

	if err := h.val.Struct(createTaskRequestDto); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.s.Create(createTaskRequestDto); err != nil {
		var appErr *app.AppError
		ok := errors.As(err, &appErr)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendSuccess(c, "Task added successfully", nil)
}

func (h *TaskHandler) SubmitTask(c *fiber.Ctx) error {
	taskId := c.Params(":id")
	if taskId == "" {
		return response.SendError(c, fiber.StatusBadRequest, "Empty Task ID", nil)
	}

	var submitTaskRequestDto entity.SubmitTaskRequestDto
	if err := c.BodyParser(&submitTaskRequestDto); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Body", err.Error())
	}

	return response.SendSuccess(c, "tasks retrieved successfully", "")
}
