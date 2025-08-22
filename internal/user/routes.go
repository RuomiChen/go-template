package user

import (
	"mvc/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func RegisterRoutes(r fiber.Router, service *Service, logger zerolog.Logger) {

	handler := NewHandler(service, logger)

	r.Post("/change-password", handler.ChangePassword) //修改密码

	// 管理员路由
	adminRoutes := r.Group("/", middleware.AdminOnly())

	adminRoutes.Get("/", handler.GetUserList)
	adminRoutes.Get("/:id", handler.GetUserDetail) // 获取新闻详情
	adminRoutes.Post("/", handler.CreateUser)
	adminRoutes.Put("/:id", handler.UpdateUser)          //全量更新
	adminRoutes.Patch("/:id", handler.PartialUpdateUser) //部分更新
	adminRoutes.Delete("/:id", handler.DeleteUser)

}
