package redis

import (
	"context"
	"time"
)

type Repository interface {
	SetKey(ctx context.Context, key string, userID string, expiration time.Duration) error
	GetKey(ctx context.Context, key string) (string, error)
	DeleteKey(ctx context.Context, key string) error
}

type repository struct {
	client *Client
}

func NewRepository(client *Client) Repository {
	return &repository{client: client}
}

func (r *repository) SetKey(ctx context.Context, key string, userID string, expiration time.Duration) error {
	return r.client.Set(ctx, key, userID, expiration)
}

func (r *repository) GetKey(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

func (r *repository) DeleteKey(ctx context.Context, key string) error {
	return r.client.Del(ctx, key)
}
