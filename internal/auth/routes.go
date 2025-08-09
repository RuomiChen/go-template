package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, jwtSecret string) {
	repo := NewRepository(db)
	service := NewService(repo, jwtSecret)
	handler := NewHandler(service, logger)

	r.Post("/login", handler.Login)
}
