package redis

import (
	"context"
	"time"
)

type Repository interface {
	SetToken(ctx context.Context, token string, userID string, expiration time.Duration) error
	GetToken(ctx context.Context, token string) (string, error)
	DeleteToken(ctx context.Context, token string) error
}

type repository struct {
	client *Client
}

func NewRepository(client *Client) Repository {
	return &repository{client: client}
}

func (r *repository) SetToken(ctx context.Context, token string, userID string, expiration time.Duration) error {
	return r.client.Set(ctx, token, userID, expiration)
}

func (r *repository) GetToken(ctx context.Context, token string) (string, error) {
	return r.client.Get(ctx, token)
}

func (r *repository) DeleteToken(ctx context.Context, token string) error {
	return r.client.Del(ctx, token)
}
