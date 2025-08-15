package tracking_event

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

func (h *Handler) RecordTrack(c *fiber.Ctx) error {
	userID := c.Locals("id").(string)

	var trackingEvent TrackingEvent
	if err := c.BodyParser(&trackingEvent); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "解析失败"})
	}

	trackingEvent.UserID = userID

	if err := h.service.Record(&trackingEvent); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.logger.Info().Interface("add trackingEvent", trackingEvent).Msg("添加一条轨迹成功！")
	return response.Success(c, nil)
}

func (h *Handler) GetUserTrackingEvents(c *fiber.Ctx) error {
	userID := c.Locals("id").(string)
	action := c.Query("action")
	if action == "" {
		return c.Status(400).JSON(fiber.Map{"error": "action 不能为空"})
	}

	// 交给 service 层处理，内部判断 action
	events, err := h.service.GetUserTrackingEventsByAction(userID, action)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return response.Success(c, events)
}
