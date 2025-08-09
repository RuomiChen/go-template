package user

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	service      Service
	logger       zerolog.Logger
	redisService redis.Service
}

func NewHandler(service Service, logger zerolog.Logger, redisService redis.Service) *Handler {
	return &Handler{service: service, logger: logger, redisService: redisService}
}

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	userID := c.Locals("user_id") // 从中间件取
	h.logger.Info().Msgf("User %v is fetching users", userID)
	users, err := h.service.GetUsers()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.logger.Info().Interface("users", users).Msg("get users success")

	return c.JSON(users)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "解析失败"})
	}
	if err := h.service.AddUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(user)
}
