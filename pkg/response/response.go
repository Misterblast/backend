package response

import (
	"github.com/gofiber/fiber/v2"
)

type PaginateResponse struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Data  any   `json:"data"`
}
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SendResponse(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Message: message,
		Data:    data,
	})
}

func SendSuccess(c *fiber.Ctx, message string, data interface{}) error {
	return SendResponse(c, fiber.StatusOK, message, data)
}

func SendError(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return SendResponse(c, statusCode, message, data)
}
