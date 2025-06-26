package handler

import (
	"errors"
	"strconv"

	"github.com/ghulammuzz/misterblast/internal/task/entity"
	"github.com/ghulammuzz/misterblast/internal/task/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type TaskSubmissionHandler struct {
	svc svc.TaskSubmissionService
	val *validator.Validate
}

func NewTaskSubmissionHandler(s svc.TaskSubmissionService, val *validator.Validate) *TaskSubmissionHandler {
	return &TaskSubmissionHandler{s, val}
}

func (h *TaskSubmissionHandler) Router(r fiber.Router) {
	r.Post("/submit-task/:id", m.R100(), m.JWTProtected(), h.SubmitTask)
	r.Put("/submission/:submissionId/score", m.R100(), h.ScoreSubmission)
	r.Get("/my-submissions", m.R100(), m.JWTProtected(), h.ListMySubmissions)
	r.Get("/task-submissions/:taskId", m.R100(), m.JWTProtected(), h.ListTaskSubmissions)
	r.Get("/submission/:submissionId", m.R100(), h.GetSubmissionDetail)
}

func (h *TaskSubmissionHandler) SubmitTask(c *fiber.Ctx) error {
	taskId, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Task ID", nil)
	}

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userId := int(claims["user_id"].(float64))

	var dto entity.SubmitTaskRequestDto
	if err := c.BodyParser(&dto); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid body", err.Error())
	}

	submitDTO := entity.SubmitTaskRequestDto{
		Answer: c.FormValue("answer"),
	}

	if form, err := c.MultipartForm(); err == nil {
		if files := form.File["attached_url"]; len(files) > 0 {
			file := files[0]
			const maxFileSize = 3 * 1024 * 1024

			if file.Size > maxFileSize {
				log.Error("File size exceeds 3MB limit", "fileSize", file.Size)
				return response.SendError(c, fiber.StatusBadRequest, "File size exceeds 3MB limit", nil)
			}

			submitDTO.AttachedURL = file
		}
	}

	if err := h.val.Struct(submitDTO); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.svc.SubmitTask(taskId, int64(userId), submitDTO); err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, err.Error())
	}

	return response.SendSuccess(c, "Task submitted successfully", nil)
}

func (h *TaskSubmissionHandler) ScoreSubmission(c *fiber.Ctx) error {
	submissionId, err := strconv.ParseInt(c.Params("submissionId"), 10, 64)
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Submission ID", nil)
	}

	// userToken := c.Locals("user").(*jwt.Token)

	// claims, ok := userToken.Claims.(jwt.MapClaims)
	// if !ok || !userToken.Valid {
	// 	log.Error("Invalid token")
	// 	return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	// }

	// userId := int(claims["user_id"].(float64))

	var dto entity.ScoreSubmissionRequestDto
	if err := c.BodyParser(&dto); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid body", err.Error())
	}

	if err := h.val.Struct(dto); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.svc.GiveScore(submissionId, dto); err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Score submitted successfully", nil)
}

func (h *TaskSubmissionHandler) ListMySubmissions(c *fiber.Ctx) error {
	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userId := int(claims["user_id"].(float64))

	filter := map[string]string{
		"page":  c.Query("page", "1"),
		"limit": c.Query("limit", "10"),
		"type":  c.Query("type", ""),
	}

	result, err := h.svc.GetSubmissionsByUser(filter, int64(userId))
	if err != nil {
		log.Error("Error retrieving submissions: %v", err)
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Submissions retrieved", result)
}

func (h *TaskSubmissionHandler) ListTaskSubmissions(c *fiber.Ctx) error {
	taskId, err := strconv.ParseInt(c.Params("taskId"), 10, 64)
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Task ID", nil)
	}

	filter := map[string]string{
		"page":  c.Query("page", "1"),
		"limit": c.Query("limit", "10"),
		"type":  c.Query("type", ""),
	}

	result, err := h.svc.GetSubmissionsByTask(filter, taskId)
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Submissions retrieved", result)
}

func (h *TaskSubmissionHandler) GetSubmissionDetail(c *fiber.Ctx) error {
	submissionId, err := strconv.ParseInt(c.Params("submissionId"), 10, 64)
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid Submission ID", nil)
	}

	result, err := h.svc.GetSubmissionDetailById(submissionId)
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Submission detail retrieved", result)
}
