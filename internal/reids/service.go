package redis

import (
	"context"
	"time"
)

type Service interface {
	SaveToken(ctx context.Context, token string, userID string, expiration time.Duration) error
	ValidateToken(ctx context.Context, token string) (string, error)
	RemoveToken(ctx context.Context, token string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) SaveToken(ctx context.Context, token string, userID string, expiration time.Duration) error {
	return s.repo.SetToken(ctx, token, userID, expiration)
}

func (s *service) ValidateToken(ctx context.Context, token string) (string, error) {
	return s.repo.GetToken(ctx, token)
}

func (s *service) RemoveToken(ctx context.Context, token string) error {
	return s.repo.DeleteToken(ctx, token)
}
