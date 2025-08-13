package tag

import (
	"errors"
	"mvc/pkg/response"
	"path/filepath"

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

func (h *Handler) GetTagList(c *fiber.Ctx) error {
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
	tag, total, err := h.service.GetTagList(page, pageSize)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get tag list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. 统一返回分页结构
	return response.Success(c, fiber.Map{
		"list":     tag,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
func (h *Handler) GetTagDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing tag id",
		})
	}

	tag, err := h.service.GetTagDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "tag not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get tag detail",
		})
	}

	return c.JSON(tag)
}
func (h *Handler) CreateTag(c *fiber.Ctx) error {
	var tag Tag
	if err := c.BodyParser(&tag); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "解析失败"})
	}
	if err := h.service.AddTag(&tag); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	h.logger.Info().Interface("add tag", tag).Msg("添加一条新新闻成功！")
	return response.Success(c, nil)
}
func (h *Handler) UpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing tag id"})
	}

	var req Tag
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	updated, err := h.service.UpdateTag(id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "tag not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updated)
}
func (h *Handler) PartialUpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing tag id"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	updated, err := h.service.PartialUpdateTag(id, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "tag not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	h.logger.Info().Interface("updated", updated).Msg("update success!")
	return response.Success(c, updated)
}
func (h *Handler) DeleteTag(c *fiber.Ctx) error {
	// 1. 从 URL 参数读取 ID
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "缺少 ID")
	}

	// 2. 调用 service 删除
	if err := h.service.DeleteTag(id); err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("删除新闻失败")
		return response.Error(c, 500, err.Error())
	}

	// 3. 成功返回
	h.logger.Info().Str("id", id).Msg("删除新闻成功")
	return response.Success(c, nil)
}
func (h *Handler) UploadImage(c *fiber.Ctx) error {
	saveDir := filepath.Join("uploads", "tag") // 统一物理目录
	imagePath, err := h.service.UploadImage(c, saveDir)
	if err != nil {
		return response.Error(c, 400, err.Error())
	}

	// 返回图片访问URL
	url := "/" + filepath.ToSlash(filepath.Join(saveDir, filepath.Base(imagePath)))
	h.logger.Info().Str("tag image url", url).Msg("upload tag image success!")
	return response.Success(c, url)
}
