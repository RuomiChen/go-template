package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Logger   zerolog.Logger // 直接在配置中暴露日志实例
}
type JWTConfig struct {
	Secret string
}
type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

// Redis 配置
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ServerConfig struct {
	Port int
}

func LoadConfig() *Config {
	v := viper.New()
	v.SetConfigName("config") // 文件名（不带扩展名）
	v.SetConfigType("yaml")
	v.AddConfigPath("./config") // 配置文件路径

	// 默认值（可选）
	v.SetDefault("server.port", 3000)
	v.SetDefault("jwt.secret", "mysecret")
	v.SetDefault("redis.addr", "localhost:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		logger.Fatal().Err(err).Msg("读取配置文件失败")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		logger.Fatal().Err(err).Msg("解析配置文件失败")
	}

	// 初始化全局日志（按需设置格式）
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger()

	cfg.Logger = logger

	return &cfg
}
