package routes

import (
	"mvc/internal/auth"
	"mvc/internal/middleware"
	"mvc/internal/news"
	"mvc/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, logger zerolog.Logger, jwtSecret string) {
	api := app.Group("/api")

	// v1 API
	v1 := api.Group("/v1")
	userGroup := v1.Group("/users", middleware.AuthMiddleware(jwtSecret))
	user.RegisterRoutes(userGroup, db, logger)

	// 未来新增模块，只需要这样
	// postGroup := v1.Group("/posts")
	// post.RegisterRoutes(postGroup, db, logger)
	newsGroup := v1.Group("/news")
	news.RegisterRoutes(newsGroup, db, logger)

	authGroup := v1.Group("/auth")
	auth.RegisterRoutes(authGroup, db, logger, jwtSecret)
}
