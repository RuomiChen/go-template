package news

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, logger)

	r.Get("/", handler.GetNewsList)
	r.Get("/:id", handler.GetNewsDetail) // 获取新闻详情
	r.Post("/", handler.CreateNews)
	r.Put("/:id", handler.UpdateNews)          //全量更新
	r.Patch("/:id", handler.PartialUpdateNews) //部分更新
	r.Delete("/:id", handler.DeleteNews)
}
