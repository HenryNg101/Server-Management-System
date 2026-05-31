package app

import (
	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/auth"
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	"github.com/HenryNg101/server-management-system/internal/platform/elastic"
	"github.com/HenryNg101/server-management-system/internal/platform/postgres"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	ElasticSession  *elasticsearch.Client
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
	postgres.MigratePostgres(postgresConfig)
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	elasticConfig := config.LoadElasticsearch()
	esSession := elastic.NewElasticsearchSession(elasticConfig)
	elastic.InitElasticsearch(esSession)

	//
	// Create services
	serverRepo := server.NewRepository(postgresSession)
	elasticServerRepo := server.NewServerESRepository(esSession)
	mailerConfig := config.LoadMailer()
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

	authService := auth.NewService(userRepo)

	return &Application{
		PostgresSession: postgresSession,
		ElasticSession:  esSession,
		ServerService:   &serverService,
		UserService:     &userService,
		AuthService:     &authService,
	}, nil
}
