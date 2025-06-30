package middleware

import (
	"fmt"
	"runtime/debug"

	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func Recover() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				stackTrace := string(debug.Stack())

				Error("[PANIC] Recovered: %v\nStack Trace:\n%s", r, stackTrace)

				_ = response.SendError(
					c,
					fiber.StatusInternalServerError,
					"[PANIC] Internal Server Error",
					fmt.Sprintf("Recovered from panic: %v", r),
				)
			}
		}()

		return c.Next()
	}
}
