package ws

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"mvc/internal/friend_relation"
	"mvc/internal/friend_request"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger) {
	// 挂载 friend 路由
	friend_request.RegisterRoutes(r, db, logger)
	friend_relation.RegisterRoutes(r, db, logger)
}
