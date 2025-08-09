package middleware

import (
	"context"
	"errors"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

// AuthMiddleware 校验 JWT 并结合 Redis 判断 Token 是否有效
func AuthMiddleware(logger zerolog.Logger, jwtSecret string, redisService redis.Service) fiber.Handler {
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

		tokenStr := parts[1]

		// 解析 JWT
		claims, err := utils.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			var ve *jwt.ValidationError
			if errors.As(err, &ve) {
				switch {
				case ve.Errors&jwt.ValidationErrorExpired != 0:
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "token expired",
					})
				case ve.Errors&jwt.ValidationErrorSignatureInvalid != 0:
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "invalid token signature",
					})
				default:
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"error": "invalid token",
					})
				}
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// --- 新增 Redis 校验部分 ---
		ctx := context.Background()
		userID, err := redisService.ValidateToken(ctx, tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "failed to validate token (redis error)",
			})
		}
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token invalid or logged out",
			})
		}
		// ---------------------------

		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}
