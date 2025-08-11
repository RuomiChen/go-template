package news

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type Handler struct {
	service Service
	logger  zerolog.Logger
}

func NewHandler(service Service, logger zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) GetNewsList(c *fiber.Ctx) error {
	news, err := h.service.GetNewsList()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.logger.Info().Interface("news", news).Msg("get news success")

	return c.JSON(news)
}

func (h *Handler) CreateNews(c *fiber.Ctx) error {
	var news News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "解析失败"})
	}
	if err := h.service.AddNews(&news); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(news)
}
