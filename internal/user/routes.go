package user

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, logger, redisService)

	r.Get("/", handler.GetUsers)
	r.Post("/", handler.CreateUser)
}
