package news_like

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	r.Post("/like/:id", handler.ToggleLike) //上传接口
}
