package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Service interface {
	SaveKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	ValidateKey(ctx context.Context, key string) (string, error)
	RemoveKey(ctx context.Context, key string) error
}

type service struct {
	repo   Repository
	logger zerolog.Logger
}

func NewService(repo Repository, logger zerolog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) SaveKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return s.repo.SetKey(ctx, key, value, expiration)
}

func (s *service) ValidateKey(ctx context.Context, key string) (string, error) {
	s.logger.Info().Str("key", key).Msg("start valid")
	key, err := s.repo.GetKey(ctx, key)
	if err == redis.Nil {
		// key 不存在或过期，正常情况，不算错误
		return "", nil
	}
	if err != nil {
		// 其他redis异常
		return "", err
	}
	return key, nil
}

func (s *service) RemoveKey(ctx context.Context, key string) error {
	return s.repo.DeleteKey(ctx, key)
}
