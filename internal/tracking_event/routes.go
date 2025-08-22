package tracking_event

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RegisterRoutes(r fiber.Router, service Service, logger zerolog.Logger, redisService redis.Service) {
	handler := NewHandler(service, logger)

	r.Get("/record", handler.GetUserTrackingEvents) //获取记录
	r.Post("/record", handler.RecordTrack)          //记录轨迹

}
