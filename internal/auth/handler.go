package auth

import (
	"mvc/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	service *Service
	logger  zerolog.Logger
}

func NewHandler(service *Service, logger zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	token, err := h.service.Login(c, req.Username, req.Password)
	h.logger.Info().Interface("token", token).Msg("login success")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.JSON(fiber.Map{"token": token})
}

// 注册接口
func (h *Handler) Register(c *fiber.Ctx) error {
	// 1. 定义请求体结构
	type RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error().Msg("invalid body")
		return response.Error(c, fiber.StatusBadRequest, "invalid body")
	}

	if req.Username == "" || req.Password == "" {
		h.logger.Error().Msg("auth username or password empty!")
		return response.Error(c, fiber.StatusBadRequest, "uth username or password empty")
	}

	// 2. 调用 service 注册
	if err := h.service.Register(c, req.Username, req.Password); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	h.logger.Info().Msg("auth register success!")
	// 3. 注册成功返回
	return response.Success(c, nil)
}
