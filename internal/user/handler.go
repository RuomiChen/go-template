package user

import (
	"errors"
	"mvc/internal/redis"
	"mvc/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Handler struct {
	service      *Service
	logger       zerolog.Logger
	redisService redis.Service
}

func NewHandler(service *Service, logger zerolog.Logger, redisService redis.Service) *Handler {
	return &Handler{service: service, logger: logger, redisService: redisService}
}

func (h *Handler) GetUserList(c *fiber.Ctx) error {
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
	user, total, err := h.service.GetUserList(page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get user list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. 统一返回分页结构
	return response.Success(c, fiber.Map{
		"list":     user,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
func (h *Handler) GetUserDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing user id",
		})
	}

	user, err := h.service.GetUserDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get user detail",
		})
	}

	return c.JSON(user)
}
func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "parse error"})
	}
	if err := h.service.AddUser(&user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.logger.Info().Interface("add user", user).Msg("add user success!")
	return response.Success(c, nil)
}
func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing user id"})
	}

	var req User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	updated, err := h.service.UpdateUser(id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updated)
}
func (h *Handler) PartialUpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing user id"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	updated, err := h.service.PartialUpdateUser(id, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "user not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	h.logger.Info().Interface("updated", updated).Msg("update success!")
	return response.Success(c, updated)
}
func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	// 1. 从 URL 参数读取 ID
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "缺少 ID")
	}

	// 2. 调用 service 删除
	if err := h.service.DeleteUser(id); err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("delete user error")
		return response.Error(c, 500, err.Error())
	}

	// 3. 成功返回
	h.logger.Info().Str("id", id).Msg("delete user success!")
	return response.Success(c, nil)
}

// 修改密码
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("id").(string)

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, 400, "请求参数错误")
	}

	if err := h.service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return response.Error(c, 400, err.Error())
	}
	return response.Success(c, "修改密码成功")
}
