package comment

import (
	"strconv"

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

func (h *Handler) Create(c *fiber.Ctx) error {
	var req struct {
		NewsId   uint64  `json:"news_id"`
		ParentID *uint64 `json:"parent_id"`
		Content  string  `json:"content"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// 从中间件拿用户ID
	userIDStr := c.Locals("id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	if err := h.service.CreateComment(userID, req.NewsId, req.ParentID, req.Content); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "comment created"})
}

func (h *Handler) List(c *fiber.Ctx) error {
	newsId, _ := strconv.ParseUint(c.Query("post_id"), 10, 64)
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	comments, err := h.service.GetComments(newsId, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(comments)
}
