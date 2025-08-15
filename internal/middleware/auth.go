package middleware

import (
	"context"
	"fmt"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func AuthMiddleware(logger zerolog.Logger, jwtSecret string, redisService redis.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing Authorization header"})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid Authorization format"})
		}

		tokenStr := parts[1]

		// 解析 JWT
		claims, err := utils.ParseToken(tokenStr, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		// Redis 校验
		ctx := context.Background()
		storedID, err := redisService.ValidateKey(ctx, tokenStr)
		if err != nil || storedID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token invalid or expired"})
		}

		idFloat, ok := claims["id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid id type"})
		}
		roleFloat, ok := claims["role"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid role type"})
		}

		// 存上下文
		c.Locals("id", fmt.Sprint(uint64(idFloat)))
		c.Locals("role", fmt.Sprint(uint64(roleFloat)))

		return c.Next()
	}
}
func AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(int)
		if !ok || role != 1 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
		}
		return c.Next()
	}
}
