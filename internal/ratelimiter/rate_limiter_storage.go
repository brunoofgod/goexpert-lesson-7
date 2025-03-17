package ratelimiter

import "time"

type RateLimiterStorage interface {
	IncrementRequestCount(key string) (int, error)
	GetRequestCount(key string) (int, error)
	SetExpiration(key string, duration time.Duration) error
	BlockKey(key string, duration time.Duration) error
	IsBlocked(key string) (bool, error)
}
