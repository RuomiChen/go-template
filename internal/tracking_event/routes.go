package tracking_event

import (
	"mvc/internal/news"
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service) {
	repo := NewRepository(db)

	newsRepo := news.NewRepository(db)
	newsService := news.NewService(newsRepo, redisService)

	service := NewService(repo, newsService, redisService)
	handler := NewHandler(service, logger)

	r.Get("/record", handler.GetUserTrackingEvents) //获取记录
	r.Post("/record", handler.RecordTrack)          //记录轨迹

}
