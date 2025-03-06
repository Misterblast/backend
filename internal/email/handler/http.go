package handler

import (
	"github.com/ghulammuzz/misterblast/internal/email/entity"
	"github.com/ghulammuzz/misterblast/internal/email/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	emailService svc.EmailService
	val          *validator.Validate
}

func NewEmailHandler(emailService svc.EmailService, val *validator.Validate) *EmailHandler {
	return &EmailHandler{emailService, val}
}

func (h *EmailHandler) Router(r fiber.Router) {
	r.Post("/activation/send-otp", h.SendOTPActivation)
	r.Post("/activation/check-otp", h.CheckOTPHandler)
}

func (h *EmailHandler) SendOTPActivation(c *fiber.Ctx) error {

	var SendOTP entity.SendOTP

	if err := c.BodyParser(&SendOTP); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.val.Struct(SendOTP); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.emailService.SendOTP(SendOTP.Email); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "OTP successfully sent to your email", nil)
}

func (h *EmailHandler) CheckOTPHandler(c *fiber.Ctx) error {

	var checkOTP entity.CheckOTP

	if err := c.BodyParser(&checkOTP); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.val.Struct(checkOTP); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.emailService.Validate(checkOTP.ID, checkOTP.OTP); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Valid", nil)
}
