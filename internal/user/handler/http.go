package handler

import (
	"time"

	"github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/internal/user/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
	m "github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type UserHandler struct {
	userService svc.UserService
	val         *validator.Validate
}

func NewUserHandler(userService svc.UserService, val *validator.Validate) *UserHandler {
	return &UserHandler{userService, val}
}

func (h *UserHandler) Router(r fiber.Router) {
	r.Post("/register", m.R100(), h.RegisterHandler)
	r.Post("/admin-check", m.R100(), h.RegisterAdminHandler)
	r.Post("/login", m.R100(), h.LoginHandler)
	r.Get("/users", m.R100(), h.ListUsersHandler)
	r.Get("/users/:id", m.R100(), h.DetailUserHandler)
	r.Delete("/users/:id", m.R100(), h.DeleteUserHandler)
	r.Put("/users/:id", m.R100(), h.EditUserHandler)
	r.Get("/me", m.JWTProtected(), m.R100(), h.MeUserHandler)
	r.Put("/reset-password", m.R100(), h.ChangePasswordHandler)
	r.Get("/summary", m.JWTProtected(), m.R100(), h.SummaryUserHandler)
	r.Put("/users/:id/password", m.R100(), h.UpdatePasswordHandler)
}

func (h *UserHandler) RegisterHandler(c *fiber.Ctx) error {

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if name == "" || email == "" || password == "" {
		return response.SendError(c, fiber.StatusBadRequest, "Missing required fields", nil)
	}

	user := entity.RegisterDTO{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if form, err := c.MultipartForm(); err == nil {
		if files := form.File["img"]; len(files) > 0 {
			file := files[0]
			const maxFileSize = 3 * 1024 * 1024

			if file.Size > maxFileSize {
				return response.SendError(c, fiber.StatusBadRequest, "File size exceeds 3MB limit", nil)
			}

			user.Img = file
		}
	}

	if err := h.val.Struct(user); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.userService.Register(user); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "User registered successfully", nil)
}

func (h *UserHandler) RegisterAdminHandler(c *fiber.Ctx) error {
	var admin entity.RegisterAdmin

	if err := c.BodyParser(&admin); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.val.Struct(admin); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.userService.RegisterAdmin(admin); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Admins registered successfully", nil)
}

func (h *UserHandler) LoginHandler(c *fiber.Ctx) error {
	var user entity.UserLogin

	if err := c.BodyParser(&user); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.val.Struct(user); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	userData, token, err := h.userService.Login(user)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * 60 * 60 * 1000),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return response.SendSuccess(c, "Login successful", userData)
}

func (h *UserHandler) ListUsersHandler(c *fiber.Ctx) error {

	filter := map[string]string{}
	if c.Query("search") != "" {
		filter["search"] = c.Query("search")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	users, err := h.userService.ListUser(filter, page, limit)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Users retrieved successfully", users)
}

func (h *UserHandler) DetailUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	user, err := h.userService.DetailUser(int32(id))
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "User retrieved successfully", user)
}

func (h *UserHandler) EditUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	user := entity.EditDTO{
		Name:  c.FormValue("name"),
		Email: c.FormValue("email"),
	}

	if form, err := c.MultipartForm(); err == nil {
		if files := form.File["img"]; len(files) > 0 {
			file := files[0]
			const maxFileSize = 3 * 1024 * 1024

			if file.Size > maxFileSize {
				return response.SendError(c, fiber.StatusBadRequest, "File size exceeds 3MB limit", nil)
			}

			user.Img = file
		}
	}

	if err := h.val.Struct(user); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.userService.EditUser(int32(id), user); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "User updated successfully", nil)
}

func (h *UserHandler) MeUserHandler(c *fiber.Ctx) error {

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	user, err := h.userService.AuthUser(int32(userID))
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "User retrieved successfully", user)
}

func (h *UserHandler) DeleteUserHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "invalid id", nil)
	}

	if err := h.userService.DeleteUser(int32(id)); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "user deleted successfully", nil)
}

func (h *UserHandler) ChangePasswordHandler(c *fiber.Ctx) error {
	var changePassword entity.ChangePassword

	if err := c.BodyParser(&changePassword); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.val.Struct(changePassword); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.userService.ChangePassword(changePassword.Token, changePassword.Password); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Change Password registered successfully", nil)
}

func (h *UserHandler) SummaryUserHandler(c *fiber.Ctx) error {
	filter := map[string]string{}

	if c.Query("lesson_id") != "" {
		filter["lesson_id"] = c.Query("lesson_id")
	}

	userToken := c.Locals("user").(*jwt.Token)

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok || !userToken.Valid {
		log.Error("Invalid token")
		return response.SendError(c, fiber.StatusUnauthorized, "Invalid token", nil)
	}

	userID := int(claims["user_id"].(float64))

	summary, err := h.userService.SummaryUser(int32(userID), filter)
	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}
	return response.SendSuccess(c, "User summary retrieved successfully", summary)
}

func (h *UserHandler) UpdatePasswordHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid user ID", nil)
	}

	var dto entity.EditPasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid JSON body", nil)
	}

	if err := h.val.Struct(dto); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
		return response.SendError(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	if err := h.userService.UpdatePassword(int32(id), dto.Password); err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}
		return response.SendError(c, appErr.Code, appErr.Message, nil)
	}

	return response.SendSuccess(c, "Password updated successfully", nil)
}
