package handler

import (
	"time"

	"github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/internal/user/svc"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
	"github.com/ghulammuzz/misterblast/pkg/middleware"
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
	r.Post("/register", h.RegisterHandler)
	r.Post("/admin-check", h.RegisterAdminHandler)
	r.Post("/login", h.LoginHandler)
	r.Get("/users", h.ListUsersHandler)
	r.Get("/users/:id", h.DetailUserHandler)
	r.Delete("/users/:id", h.DeleteUserHandler)
	r.Put("/users/:id", h.EditUserHandler)
	r.Get("/me", middleware.JWTProtected(), h.MeUserHandler)
}

func (h *UserHandler) RegisterHandler(c *fiber.Ctx) error {
	var user entity.Register

	if err := c.BodyParser(&user); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
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

	var user entity.EditUser
	if err := c.BodyParser(&user); err != nil {
		return response.SendError(c, fiber.StatusBadRequest, "Invalid request body", nil)
	}

	if err := h.val.Struct(user); err != nil {
		validationErrors := app.ValidationErrorResponse(err)
		log.Error("Validation failed: %v", validationErrors)
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
