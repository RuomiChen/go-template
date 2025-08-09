package friend

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {
	service := NewFriendService(db)
	handler := NewFriendHandler(service, logger)

	r.Get("/friend", websocket.New(handler.Handle))
}
