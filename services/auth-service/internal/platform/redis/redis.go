package redis

import (
	"fmt"

	"github.com/HenryNg101/auth-service/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewPostgresSession(config *config.RedisConfig) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
	})
	return redisClient
}
