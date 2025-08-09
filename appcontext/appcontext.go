package appcontext

import (
	"mvc/internal/redis" // 你的redis包路径

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type AppContext struct {
	DB           *gorm.DB
	Logger       zerolog.Logger
	RedisService redis.Service
	JWTSecret    string
}

func NewAppContext(db *gorm.DB, logger zerolog.Logger, redisService redis.Service, jwtSecret string) *AppContext {
	return &AppContext{
		DB:           db,
		Logger:       logger,
		RedisService: redisService,
		JWTSecret:    jwtSecret,
	}
}
