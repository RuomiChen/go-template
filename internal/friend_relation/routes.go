package friend_relation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service, logger)

	r.Get("/friend_relation", websocket.New(handler.Handle))

}
