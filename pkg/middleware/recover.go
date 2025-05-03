package middleware

import (
	"fmt"

	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				Error("Recovered from panic: %v", err)
				response.SendError(c, fiber.StatusInternalServerError, "[PANIC] Internal Server Error", fmt.Sprintf("%v", err))
			}
		}()

		return c.Next()
	}
}
