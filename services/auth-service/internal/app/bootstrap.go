package app

import (
	"errors"

	"github.com/HenryNg101/auth-service/internal/auth"
	"github.com/HenryNg101/auth-service/internal/config"
	"github.com/HenryNg101/auth-service/internal/platform/postgres"
	redisServer "github.com/HenryNg101/auth-service/internal/platform/redis"
	"github.com/HenryNg101/auth-service/internal/user"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	RedisSession    *redis.Client
	AuthService     auth.Service
}

func NewApp() (*Application, error) {
	// Load config
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	redisConfig := config.LoadRedis()
	if len(redisConfig.Password) == 0 {
		return nil, errors.New("Password for Redis connection is not set. You have to set it in .env file in root folder using REDIS_PASSWORD variable")
	}
	postgresSession := postgres.NewPostgresSession(postgresConfig)
	redisSession := redisServer.NewRedisSession(redisConfig)

	userRepo := user.NewRepository(postgresSession)
	authRepo := auth.NewRepository(redisSession)

	authService := auth.NewService(userRepo, authRepo)

	return &Application{
		PostgresSession: postgresSession,
		RedisSession:    redisSession,
		AuthService:     authService,
	}, nil
}
