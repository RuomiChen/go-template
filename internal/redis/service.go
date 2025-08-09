package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Service interface {
	SaveToken(ctx context.Context, token string, userID string, expiration time.Duration) error
	ValidateToken(ctx context.Context, token string) (string, error)
	RemoveToken(ctx context.Context, token string) error
}

type service struct {
	repo   Repository
	logger zerolog.Logger
}

func NewService(repo Repository, logger zerolog.Logger) Service {
	return &service{repo: repo, logger: logger}
}

func (s *service) SaveToken(ctx context.Context, token string, userID string, expiration time.Duration) error {
	return s.repo.SetToken(ctx, token, userID, expiration)
}

func (s *service) ValidateToken(ctx context.Context, token string) (string, error) {
	s.logger.Info().Str("valid token:", token).Msg("validate token")
	userID, err := s.repo.GetToken(ctx, token)
	s.logger.Info().Str("userID:", userID).Msg("validate token")
	if err == redis.Nil {
		// key 不存在或过期，正常情况，不算错误
		return "", nil
	}
	if err != nil {
		// 其他redis异常
		return "", err
	}
	return userID, nil
}

func (s *service) RemoveToken(ctx context.Context, token string) error {
	return s.repo.DeleteToken(ctx, token)
}
