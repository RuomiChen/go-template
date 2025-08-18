package tracking_event

import (
	"mvc/internal/news"
	"mvc/internal/news_like"
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service) {
	repo := NewRepository(db)

	news_like_repo := news_like.NewRepository(db)
	news_like_service := news_like.NewService(news_like_repo)

	newsRepo := news.NewRepository(db)
	newsService := news.NewService(newsRepo, redisService, news_like_service)

	service := NewService(repo, newsService, redisService)
	handler := NewHandler(service, logger)

	r.Get("/record", handler.GetUserTrackingEvents) //获取记录
	r.Post("/record", handler.RecordTrack)          //记录轨迹

}
