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
	var auth Auth
	if err := c.BodyParser(&auth); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request")
	}

	token, err := h.service.Login(c, auth.Username, auth.Password)
	h.logger.Info().Interface("token", token).Msg("login success")
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, "invalid credentials")
	}

	return response.Success(c, token)
}

// 注册接口
func (h *Handler) Register(c *fiber.Ctx) error {
	var auth Auth

	if err := c.BodyParser(&auth); err != nil {
		h.logger.Error().Msg("invalid body")
		return response.Error(c, fiber.StatusBadRequest, "invalid body")
	}
	h.logger.Info().Interface("auth", auth).Msg("123")
	if auth.Username == "" || auth.Password == "" {
		h.logger.Error().Msg("auth username or password empty!")
		return response.Error(c, fiber.StatusBadRequest, "uth username or password empty")
	}

	// 2. 调用 service 注册
	if err := h.service.Register(c, auth.Username, auth.Password); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}
	h.logger.Info().Msg("auth register success!")
	// 3. 注册成功返回
	return response.Success(c, nil)
}
