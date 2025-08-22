package comment

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, logger)

	group := r.Group("/comments")
	group.Post("/", handler.Create)
	group.Get("/", handler.List)
}
