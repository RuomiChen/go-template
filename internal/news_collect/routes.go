package news_collect

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RegisterRoutes(r fiber.Router, service Service, logger zerolog.Logger, redisService redis.Service) {
	handler := NewHandler(service)

	r.Post("/collect/:id", handler.ToggleLike) //上传接口
	r.Get("/collects", handler.GetUserCollects)
}
