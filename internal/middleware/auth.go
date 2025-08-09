package middleware

import (
	"errors"
	"mvc/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing Authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid Authorization format",
			})
		}

		claims, err := utils.ParseToken(parts[1], jwtSecret)
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				switch {
				case ve.Errors&jwt.ValidationErrorExpired != 0:
					// Token 过期
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "token expired",
					})
				case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
					// 签名不合法
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "invalid token signature",
					})
				default:
					// 其他验证错误
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "invalid token",
					})
				}
			}

			// 其他非 jwt.ValidationError 的错误
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}
