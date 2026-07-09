package auth

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisServerRepository interface {
	StoreToken(ctx context.Context, key string, value []byte, ttl time.Duration) error
	GetUserInfo(ctx context.Context, key string) (string, error)
	DeleteOldRefreshToken(ctx context.Context, key string) error
}

type redisServerRepository struct {
	redisClient *redis.Client
}

func NewRepository(client *redis.Client) RedisServerRepository {
	return &redisServerRepository{redisClient: client}
}

func (r *redisServerRepository) StoreToken(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.redisClient.Set(ctx, key, value, ttl).Err()
}

func (r *redisServerRepository) GetUserInfo(ctx context.Context, key string) (string, error) {
	return r.redisClient.Get(ctx, key).Result()
}

func (r *redisServerRepository) DeleteOldRefreshToken(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
