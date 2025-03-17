package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiterStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRateLimiterStorage(client *redis.Client) *RedisRateLimiterStorage {
	return &RedisRateLimiterStorage{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisRateLimiterStorage) IncrementRequestCount(key string) (int, error) {
	count, err := r.client.Incr(r.ctx, key).Result()

	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *RedisRateLimiterStorage) GetRequestCount(key string) (int, error) {
	count, err := r.client.Get(r.ctx, key).Int()

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RedisRateLimiterStorage) SetExpiration(key string, duration time.Duration) error {
	return r.client.Expire(r.ctx, key, duration).Err()
}

func (r *RedisRateLimiterStorage) BlockKey(key string, duration time.Duration) error {
	return r.client.Expire(r.ctx, key, duration).Err()
}

func (r *RedisRateLimiterStorage) IsBlocked(key string) (bool, error) {
	exists, err := r.client.Exists(r.ctx, "block:"+key).Result()
	return exists > 0, err
}
