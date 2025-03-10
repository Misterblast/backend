package middleware

import (
	"github.com/ghulammuzz/misterblast/pkg/jwt"
	"github.com/ghulammuzz/misterblast/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return response.SendError(c, 401, "Unauthorized", "token not found")
		}

		token, claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			return response.SendError(c, 401, "Unauthorized", err.Error())
		}

		c.Locals("user", token)
		c.Locals("claims", claims)

		return c.Next()
	}
}
