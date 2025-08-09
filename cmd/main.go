package main

import (
	"context"
	"fmt"
	"mvc/config"
	redis "mvc/internal/reids"
	"mvc/routes"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	cfg := config.LoadConfig()

	// 初始化 Redis
	redisClient := redis.NewRedisClient(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	redisRepo := redis.NewRepository(redisClient)
	redisService := redis.NewService(redisRepo)
	if !redisClient.IsConnected(context.Background()) {
		cfg.Logger.Fatal().Msg("Redis connection failed")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		cfg.Logger.Fatal().Err(err).Msg("数据库连接失败")
	}

	app := fiber.New()
	routes.Register(app, db, cfg.Logger, redisService, cfg.JWT.Secret)

	cfg.Logger.Info().Int("port", cfg.Server.Port).Msg("服务器启动中...")
	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}
