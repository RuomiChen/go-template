package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"mvc/internal/friend_relation"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {

	friend_relation.RegisterRoutes(r, db, logger)
}
