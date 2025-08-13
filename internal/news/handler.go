package news

import (
	"errors"
	"mvc/pkg/response"
	"path/filepath"
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

func (h *Handler) GetTopNews(c *fiber.Ctx) error {
	n := c.QueryInt("n", 10)
	if n <= 0 {
		n = 10
	}

	news, err := h.service.GetTopNews(n)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get news list")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return response.Success(c, news)
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
func (h *Handler) GetNewsDetail(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing news id",
		})
	}

	news, err := h.service.GetNewsDetail(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "news not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get news detail",
		})
	}

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
	h.logger.Info().Interface("add news", news).Msg("添加一条新新闻成功！")
	return response.Success(c, nil)
}
func (h *Handler) UpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing news id"})
	}

	var req News
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	updated, err := h.service.UpdateNews(id, &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "news not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(updated)
}
func (h *Handler) PartialUpdateNews(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing news id"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "invalid request body")
	}

	updated, err := h.service.PartialUpdateNews(id, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.Error(c, fiber.StatusNotFound, "news not found")
		}
		return response.Error(c, fiber.StatusInternalServerError, err.Error())
	}
	h.logger.Info().Interface("updated", updated).Msg("update success!")
	return response.Success(c, updated)
}
func (h *Handler) DeleteNews(c *fiber.Ctx) error {
	// 1. 从 URL 参数读取 ID
	id := c.Params("id")
	if id == "" {
		return response.Error(c, 400, "缺少 ID")
	}

	// 2. 调用 service 删除
	if err := h.service.DeleteNews(id); err != nil {
		h.logger.Error().Err(err).Str("id", id).Msg("删除新闻失败")
		return response.Error(c, 500, err.Error())
	}

	// 3. 成功返回
	h.logger.Info().Str("id", id).Msg("删除新闻成功")
	return response.Success(c, nil)
}
func (h *Handler) UploadImage(c *fiber.Ctx) error {
	saveDir := filepath.Join("uploads", "news") // 统一物理目录
	imagePath, err := h.service.UploadImage(c, saveDir)
	if err != nil {
		return response.Error(c, 400, err.Error())
	}

	// 返回图片访问URL
	url := "/" + filepath.ToSlash(filepath.Join(saveDir, filepath.Base(imagePath)))
	h.logger.Info().Str("news image url", url).Msg("upload news image success!")
	return response.Success(c, url)
}

func (h *Handler) GetNewsByTag(c *fiber.Ctx) error {
	tagIDStr := c.Params("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid tagId"})
	}

	limitStr := c.Query("size", "10") // 每页大小
	pageStr := c.Query("page", "1")   // 页码，从 1 开始

	limit, _ := strconv.Atoi(limitStr)
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit // 计算偏移量

	newsList, err := h.service.GetNewsByTag(uint(tagID), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})

	}
	h.logger.Info().Interface("news", newsList).Msg("succ")

	return response.Success(c, newsList)
}
