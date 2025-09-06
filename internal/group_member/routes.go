package group_member

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RegisterRoutes(r fiber.Router, service Service, logger zerolog.Logger) {
	handler := NewHandler(service, logger)

	// 创建群组
	r.Get("/join/:id", handler.JoinGroup)
	//获取当前用户关联所有群组
	r.Get("/leave", handler.LeaveGroup)
}
