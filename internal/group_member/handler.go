package group_member

import (
	"fmt"
	"mvc/pkg/response"
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

// 加入群组
func (h *Handler) JoinGroup(c *fiber.Ctx) error {

	groupID, _ := strconv.ParseUint(c.Params("id"), 10, 64)
	fmt.Print(groupID)

	userIDStr := c.Locals("id")
	if userIDStr == nil {
		return response.Error(c, fiber.StatusUnauthorized, "missing user id")
	}
	userID, _ := strconv.ParseUint(userIDStr.(string), 10, 64)

	if err := h.service.JoinGroup(groupID, userID, "member"); err != nil {
		h.logger.Error().Err(err).Msg("join group failed")
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	h.logger.Info().Uint64("group", groupID).Uint64("user", userID).Msg("user joined group")
	return response.Success(c, fiber.Map{"message": "joined group successfully"})
}

// 退出群组
func (h *Handler) LeaveGroup(c *fiber.Ctx) error {
	groupID, _ := strconv.ParseUint(c.Params("id"), 10, 64)

	userIDStr := c.Locals("id")
	if userIDStr == nil {
		return response.Error(c, fiber.StatusUnauthorized, "missing user id")
	}
	userID, _ := strconv.ParseUint(userIDStr.(string), 10, 64)

	if err := h.service.LeaveGroup(groupID, userID); err != nil {
		h.logger.Error().Err(err).Msg("leave group failed")
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	h.logger.Info().Uint64("group", groupID).Uint64("user", userID).Msg("user left group")
	return response.Success(c, fiber.Map{"message": "left group successfully"})
}

// 获取群组成员列表
func (h *Handler) GetMembers(c *fiber.Ctx) error {
	groupID, _ := strconv.ParseUint(c.Params("id"), 10, 64)

	members, err := h.service.GetGroupMembers(groupID)
	if err != nil {
		h.logger.Error().Err(err).Msg("get group members failed")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, members)
}
