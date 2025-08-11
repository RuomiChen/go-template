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
	r.Post("/", handler.CreateNews)
}
