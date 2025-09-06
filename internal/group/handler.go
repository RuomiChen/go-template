package group

import (
	"errors"
	"mvc/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
	logger  zerolog.Logger
}

func NewHandler(service Service, logger zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// 获取群组列表（分页）
func (h *Handler) GetGroupList(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	groups, total, err := h.service.GetGroupList(page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get group list")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, fiber.Map{
		"list":     groups,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// 获取群组详情
func (h *Handler) GetGroupDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	groupID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid group id")
	}

	group, err := h.service.GetGroupDetail(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "group not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	return response.Success(c, group)
}

// 创建群组（不用事务）
// 会同时创建群聊并将创建者加入群成员表，角色是 owner
func (h *Handler) CreateGroup(c *fiber.Ctx) error {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	userIDStr := c.Locals("id")
	if userIDStr == nil {
		return response.Error(c, fiber.StatusUnauthorized, "missing user id")
	}
	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id")
	}

	group, err := h.service.CreateGroupWithOwner(req.Name, userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	h.logger.Info().
		Uint64("owner", userID).
		Str("name", req.Name).
		Msg("create group success")

	return response.Success(c, group)
}

// 获取当前用户所在的所有群组（普通/管理员/群主）
func (h *Handler) GetUserGroups(c *fiber.Ctx) error {
	// 1. 从 Locals 获取当前登录用户 ID
	userIDStr := c.Locals("id")
	if userIDStr == nil {
		return response.Error(c, fiber.StatusUnauthorized, "missing user id")
	}

	userID, err := strconv.ParseUint(userIDStr.(string), 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid user id")
	}

	// 2. 调用 Service 层获取群组列表
	groups, err := h.service.GetGroupsByUser(userID)
	if err != nil {
		h.logger.Error().Err(err).Uint64("userID", userID).Msg("failed to get user groups")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	// 3. 返回列表
	return response.Success(c, groups)
}

// 更新群组
func (h *Handler) UpdateGroup(c *fiber.Ctx) error {
	groupID, _ := strconv.ParseUint(c.Params("id"), 10, 64)

	var req Group
	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	updated, err := h.service.UpdateGroup(groupID, req.Name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "group not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	h.logger.Info().Uint64("id", groupID).Str("name", req.Name).Msg("update group success")
	return response.Success(c, updated)
}

// 删除群组
func (h *Handler) DeleteGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	groupID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid group id")
	}

	if err := h.service.DeleteGroup(groupID); err != nil {
		h.logger.Error().Err(err).Uint64("id", groupID).Msg("delete group failed")
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}

	h.logger.Info().Uint64("id", groupID).Msg("delete group success")
	return response.Success(c, nil)
}
