package friend_request

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewFriendHandler(service, logger)

	r.Get("/friend_request", websocket.New(handler.Handle))
}
