package auth

import (
	"mvc/internal/admin"
	redis "mvc/internal/redis"
	"mvc/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func RegisterRoutes(r fiber.Router, db *gorm.DB, logger zerolog.Logger, redisService redis.Service, jwtSecret string) {

	adminRepo := admin.NewRepository(db)
	adminService := admin.NewService(adminRepo, jwtSecret, redisService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, jwtSecret, redisService)

	handler := NewHandler(adminService, userService)

	r.Post("/login", handler.Login)
	r.Post("/register", handler.Register)

}
