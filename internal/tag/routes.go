package tag

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service) {
	repo := NewRepository(db)
	service := NewService(repo, redisService)
	handler := NewHandler(service, logger)

	r.Post("/upload", handler.UploadImage) //上传接口

	r.Get("/", handler.GetTagList)      // 分页获取
	r.Get("/:id", handler.GetTagDetail) // 获取新闻详情
	r.Post("/", handler.CreateTag)
	r.Put("/:id", handler.UpdateTag)          //全量更新
	r.Patch("/:id", handler.PartialUpdateTag) //部分更新
	r.Delete("/:id", handler.DeleteTag)
}
