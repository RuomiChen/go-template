package news_like

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RegisterRoutes(r fiber.Router, service *Service, logger zerolog.Logger, redisService redis.Service) {
	handler := NewHandler(service)

	r.Post("/like/:id", handler.ToggleLike) //上传接口
}
