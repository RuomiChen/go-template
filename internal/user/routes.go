package user

import (
	"mvc/internal/middleware"
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service, jwtSecret string) {
	repo := NewRepository(db)
	service := NewService(repo, jwtSecret, redisService)
	handler := NewHandler(service, logger, redisService)

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
