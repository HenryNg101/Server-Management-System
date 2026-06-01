package app

import (
	"errors"

	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/auth"
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	"github.com/HenryNg101/server-management-system/internal/platform/elastic"
	"github.com/HenryNg101/server-management-system/internal/platform/postgres"
	redisServer "github.com/HenryNg101/server-management-system/internal/platform/redis"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	ElasticSession  *elasticsearch.Client
	RedisSession    *redis.Client
	ServerService   *server.Service
	UserService     *user.Service
	AuthService     *auth.Service
}

// Create a new app with services to be used
// Why I do this ? So that, when an app is created, it can use all internal functionalities without having to use other ways of communication like gRPC or HTTP
// If it's API app, you can add handlers, add API group, etc. to it. If it's a worker, it can still get access to DB query, services, etc. without it being separated entirely
// It might be tightly coupled, but it's good for now
func NewApp() (*Application, error) {
	//
	// Load configs and create initial schema/data streams/indices
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	postgres.MigratePostgres(postgresConfig)
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	elasticConfig := config.LoadElasticsearch()
	if len(elasticConfig.Password) == 0 {
		return nil, errors.New("Password for Elasticsearch connection is not set. You have to set it in .env file in root folder using ELASTIC_PASSWORD variable")
	}
	esSession := elastic.NewElasticsearchSession(elasticConfig)
	elastic.InitElasticsearch(esSession)

	redisConfig := config.LoadRedis()
	if len(redisConfig.Password) == 0 {
		return nil, errors.New("Password for Redis connection is not set. You have to set it in .env file in root folder using REDIS_PASSWORD variable")
	}
	redisSession := redisServer.NewPostgresSession(redisConfig)

	//
	// Create services
	serverRepo := server.NewRepository(postgresSession)
	elasticServerRepo := server.NewServerESRepository(esSession)
	mailerConfig := config.LoadMailer()
	if len(mailerConfig.Password) == 0 {
		return nil, errors.New("Password for email server is not set. You have to set it in .env file in root folder using MAIL_PASSWORD variable")
	}
	mailerUtility := server.NewMailer(
		mailerConfig.Server,
		mailerConfig.Port,
		mailerConfig.UserName,
		mailerConfig.Password,
		mailerConfig.FromEmail,
	)
	serverService := server.NewService(serverRepo, elasticServerRepo, mailerUtility)

	userRepo := user.NewRepository(postgresSession)
	userService := user.NewService(userRepo)

	authRedisRepo := auth.NewRepository(redisSession)
	authService := auth.NewService(userRepo, authRedisRepo)

	return &Application{
		PostgresSession: postgresSession,
		ElasticSession:  esSession,
		RedisSession:    redisSession,
		ServerService:   &serverService,
		UserService:     &userService,
		AuthService:     &authService,
	}, nil
}
