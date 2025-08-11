package news

import (
	"mvc/pkg/response"

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
	// 1. 读取分页参数，设置默认值
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	// 2. 调用 service
	news, total, err := h.service.GetNewsList(page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get news list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. 统一返回分页结构
	return response.Success(c, fiber.Map{
		"list":     news,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

func (h *Handler) CreateNews(c *fiber.Ctx) error {
	var news News
	if err := c.BodyParser(&news); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "解析失败"})
	}
	if err := h.service.AddNews(&news); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.logger.Info().Interface("add news", news).Msg("添加一条新新闻成功！")
	return response.Success(c, nil)
}
