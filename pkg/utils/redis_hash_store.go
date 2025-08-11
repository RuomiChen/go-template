package utils

import (
	"context"
	"time"

	"mvc/internal/redis"
)

type RedisHashStore struct {
	redisService redis.Service // 用接口类型
	prefix       string
	ttl          time.Duration
}

func NewRedisHashStore(redisService redis.Service, prefix string, ttl time.Duration) *RedisHashStore {
	return &RedisHashStore{
		redisService: redisService,
		prefix:       prefix,
		ttl:          ttl,
	}
}

// Get 从 redis 取值，返回值和是否存在
func (r *RedisHashStore) Get(hash string) (string, bool) {
	ctx := context.Background()
	key := r.prefix + hash
	val, err := r.redisService.ValidateToken(ctx, key) // 调用接口里的Get操作
	if err != nil || val == "" {
		return "", false
	}
	return val, true
}

// Set 设置缓存
func (r *RedisHashStore) Set(hash, path string) {
	ctx := context.Background()
	key := r.prefix + hash
	_ = r.redisService.SaveToken(ctx, key, path, r.ttl) // 调用接口里的Set操作
}
