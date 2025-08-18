package news

import (
	"mvc/internal/news_like"
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service) {
	news_like_repo := news_like.NewRepository(db)
	news_like_service := news_like.NewService(news_like_repo)

	repo := NewRepository(db)
	service := NewService(repo, redisService, news_like_service)
	handler := NewHandler(service, logger)

	r.Post("/upload", handler.UploadImage) //上传接口

	r.Get("/", handler.GetNewsList)          // 分页获取新闻详情
	r.Get("/top", handler.GetTopNews)        // 分页获取新闻详情
	r.Get("/tags/:id", handler.GetNewsByTag) // 根据标签获取新闻
	r.Get("/:id", handler.GetNewsDetail)     // 获取新闻详情

	r.Post("/", handler.CreateNews)
	r.Put("/:id", handler.UpdateNews)          //全量更新
	r.Patch("/:id", handler.PartialUpdateNews) //部分更新
	r.Delete("/:id", handler.DeleteNews)
}
