package news_collect

import (
	"mvc/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ToggleLike(c *fiber.Ctx) error {
	newsID, _ := strconv.ParseUint(c.Params("id"), 10, 64)

	// 从 Locals 获取 userID，假设是 string
	idStr := c.Locals("id").(string)
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	// 调用 service 判断当前状态
	liked, err := h.svc.ToggleCollect(c.Context(), newsID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	msg := "收藏成功"
	if !liked {
		msg = "取消收藏"
	}
	return response.Success(c, msg)
}
func (h *Handler) GetUserCollects(c *fiber.Ctx) error {
	idStr := c.Locals("id").(string)
	userID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	collects, err := h.svc.GetUserCollects(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return response.Success(c, collects)
}
