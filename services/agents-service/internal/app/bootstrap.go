package app

import (
	"errors"
	"time"

	"github.com/HenryNg101/agents-service/internal/agent"
	"github.com/HenryNg101/agents-service/internal/config"
	"github.com/HenryNg101/agents-service/internal/platform/cache"
	kafkaClient "github.com/HenryNg101/agents-service/internal/platform/kafka"
	"github.com/HenryNg101/agents-service/internal/platform/postgres"
	"github.com/HenryNg101/agents-service/internal/server"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	RedisSession    *redis.Client
	L1Cache         *cache.LocalCache
	AgentService    agent.Service
}

func NewApp() (*Application, error) {
	// Load config
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	kafkaConfig := config.LoadKafka()
	if len(kafkaConfig.Brokers) == 0 {
		return nil, errors.New("Kafka brokers not configured. You have to set it in .env file in root folder using KAFKA_BROKERS variable")
	}

	localCache := cache.NewLocalCache(5*time.Minute, 20000)

	redisConfig := config.LoadRedis()
	if len(redisConfig.Password) == 0 {
		return nil, errors.New("Password for Redis connection is not set. You have to set it in .env file in root folder using REDIS_PASSWORD variable")
	}
	redisSession := cache.NewRedisSession(redisConfig)

	agentRepo := agent.NewRepository(postgresSession)
	serverRepo := server.NewRepository(postgresSession)
	kafkaAgentProducer := agent.NewKafkaProducer(kafkaClient.NewProducer(kafkaConfig.Brokers, kafkaConfig.AgentMetricsTopic))
	agentService := agent.NewService(agentRepo, kafkaAgentProducer, serverRepo, redisSession)

	return &Application{
		PostgresSession: postgresSession,
		RedisSession:    redisSession,
		L1Cache:         localCache,
		AgentService:    agentService,
	}, nil
}
