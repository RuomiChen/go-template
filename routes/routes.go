package routes

import (
	"mvc/appcontext"
	"mvc/internal/auth"
	"mvc/internal/middleware"
	"mvc/internal/news"
	"mvc/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Register(app *fiber.App, appCtx *appcontext.AppContext) {
	api := app.Group("/api")

	// v1 API
	v1 := api.Group("/v1")
	userGroup := v1.Group("/users", middleware.AuthMiddleware(appCtx.Logger, appCtx.JWTSecret, appCtx.RedisService))
	user.RegisterRoutes(userGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService)

	// 未来新增模块，只需要这样
	// postGroup := v1.Group("/posts")
	// post.RegisterRoutes(postGroup, db, logger)
	newsGroup := v1.Group("/news")
	news.RegisterRoutes(newsGroup, appCtx.DB, appCtx.Logger)

	authGroup := v1.Group("/auth")
	auth.RegisterRoutes(authGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService, appCtx.JWTSecret)

	api.Get("/ws", websocket.New(func(c *websocket.Conn) {
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				break
			}

			appCtx.Logger.Info().Msgf("WebSocket recv: %s", string(msg))

			// 简单回显
			if err := c.WriteMessage(mt, msg); err != nil {
				break
			}
		}
	}))
}
