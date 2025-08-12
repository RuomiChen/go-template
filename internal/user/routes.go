package user

import (
	"mvc/internal/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service, jwtSecret string) {
	repo := NewRepository(db)
	service := NewService(repo, jwtSecret, redisService)
	handler := NewHandler(service, logger, redisService)

	r.Get("/", handler.GetUserList)
	r.Get("/:id", handler.GetUserDetail) // 获取新闻详情
	r.Post("/", handler.CreateUser)
	r.Put("/:id", handler.UpdateUser)          //全量更新
	r.Patch("/:id", handler.PartialUpdateUser) //部分更新
	r.Delete("/:id", handler.DeleteUser)

}
