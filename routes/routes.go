package routes

import (
	"mvc/appcontext"
	"mvc/internal/auth"
	"mvc/internal/middleware"
	"mvc/internal/news"
	"mvc/internal/user"
	"mvc/internal/ws"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, appCtx *appcontext.AppContext) {

	app.Static("/uploads", "./uploads")

	api := app.Group("/api")

	// v1 API
	v1 := api.Group("/v1")
	userGroup := v1.Group("/users", middleware.AuthMiddleware(appCtx.Logger, appCtx.JWTSecret, appCtx.RedisService))
	user.RegisterRoutes(userGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService, appCtx.JWTSecret)

	// 未来新增模块，只需要这样
	// postGroup := v1.Group("/posts")
	// post.RegisterRoutes(postGroup, db, logger)
	newsGroup := v1.Group("/news")
	news.RegisterRoutes(newsGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService)

	authGroup := v1.Group("/auth")
	auth.RegisterRoutes(authGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService, appCtx.JWTSecret)

	wsGroup := v1.Group("/ws")
	ws.RegisterRoutes(wsGroup, appCtx.DB, appCtx.Logger)
}
