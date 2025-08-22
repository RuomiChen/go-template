package routes

import (
	"mvc/appcontext"
	"mvc/internal/auth"
	"mvc/internal/friend_request"
	"mvc/internal/middleware"
	"mvc/internal/news"
	"mvc/internal/news_like"
	"mvc/internal/tag"
	"mvc/internal/tracking_event"
	"mvc/internal/user"
	"mvc/internal/ws"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App, appCtx *appcontext.AppContext) {

	app.Static("/uploads", "./uploads")

	api := app.Group("/api")

	// v1 API
	v1 := api.Group("/v1")

	newsLikeRepo := news_like.NewRepository(appCtx.DB)
	newsLikeService := news_like.NewService(newsLikeRepo)

	newsRepo := news.NewRepository(appCtx.DB)
	newsService := news.NewService(newsRepo, appCtx.RedisService, newsLikeService)

	trackingEventRepo := tracking_event.NewRepository(appCtx.DB)
	trackingEventService := tracking_event.NewService(trackingEventRepo, newsService, appCtx.RedisService)

	trackingGroup := v1.Group("/tracking", middleware.AuthMiddleware(appCtx.Logger, appCtx.JWTSecret, appCtx.RedisService))
	tracking_event.RegisterRoutes(trackingGroup, trackingEventService, appCtx.Logger, appCtx.RedisService)

	userRepo := user.NewRepository(appCtx.DB)
	userService := user.NewService(userRepo, appCtx.JWTSecret, appCtx.RedisService)

	userGroup := v1.Group("/users", middleware.AuthMiddleware(appCtx.Logger, appCtx.JWTSecret, appCtx.RedisService))

	user.RegisterRoutes(userGroup, userService, appCtx.Logger)

	newsGroup := v1.Group("/news", middleware.OptionalAuthMiddleware(appCtx.JWTSecret, appCtx.RedisService))
	news.RegisterRoutes(newsGroup, newsService, appCtx.Logger, appCtx.RedisService)

	newsLikeGroup := v1.Group("/news_like", middleware.AuthMiddleware(appCtx.Logger, appCtx.JWTSecret, appCtx.RedisService))
	news_like.RegisterRoutes(newsLikeGroup, newsLikeService, appCtx.Logger, appCtx.RedisService)

	tagGroup := v1.Group("/tag")
	tag.RegisterRoutes(tagGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService)

	authGroup := v1.Group("/auth")
	auth.RegisterRoutes(authGroup, appCtx.DB, appCtx.Logger, appCtx.RedisService, appCtx.JWTSecret)

	friendRequestGroup := v1.Group("/friend_request", middleware.AuthMiddleware(appCtx.Logger, appCtx.JWTSecret, appCtx.RedisService))
	friend_request.RegisterRoutes(friendRequestGroup, appCtx.DB, appCtx.Logger)

	wsGroup := v1.Group("/ws")
	ws.RegisterRoutes(wsGroup, appCtx.DB, appCtx.Logger)

}
