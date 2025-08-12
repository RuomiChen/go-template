package main

import (
	"fmt"
	"mvc/appcontext"
	"mvc/config"
	"mvc/internal/redis"
	"mvc/pkg/utils"
	"mvc/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	cfg := config.LoadConfig()

	log := utils.NewLogger(cfg.Logger.Base)

	// 初始化 Redis
	redisClient := redis.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	redisRepo := redis.NewRepository(redisClient)
	redisService := redis.NewService(redisRepo, log)

	// if !redisClient.IsConnected(context.Background()) {
	// 	log.Fatal().Msg("Redis connection failed")
	// }
	log.Info().Msg("Redis初始化成功")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Name,
	)

	db := appcontext.InitDB(dsn)

	//使用上下文的方式 避免 参数地狱
	/**
	routes.Register(app, ctx) 这样传一个上下文对象，就不会有一堆参数。
	优点：简化函数签名。
	缺点：依赖变成隐式（调某个 handler 时，不太容易看出它依赖了哪些）。
	*/
	appCtx := appcontext.NewAppContext(db, log, redisService, cfg.JWT.Secret)

	app := fiber.New()

	routes.Register(app, appCtx)

	// 直接在这里写 websocket 路由

	log.Info().Int("port", cfg.Server.Port).Msg("服务器启动中...")

	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
