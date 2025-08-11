package auth

import (
	redis "mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service, jwtSecret string) {
	repo := NewRepository(db)

	service := NewService(repo, jwtSecret, redisService)
	handler := NewHandler(service, logger)

	r.Post("/login", handler.Login)
	r.Post("/register", handler.Register)
}
