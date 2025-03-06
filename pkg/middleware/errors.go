package middleware

import (
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware(c *fiber.Ctx) error {
	err := c.Next()

	if err != nil {
		appErr, ok := err.(*app.AppError)
		if !ok {
			appErr = app.ErrInternal
		}

		return response.SendError(c, appErr.Code, appErr.Message, err.Error())
	}

	return nil
}
