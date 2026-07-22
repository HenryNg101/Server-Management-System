package app

import (
	"errors"

	"github.com/HenryNg101/server-service/internal/config"
	"github.com/HenryNg101/server-service/internal/platform/postgres"
	"github.com/HenryNg101/server-service/internal/server"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	ServerService   server.Service
}

func NewApp() (*Application, error) {
	// Load config
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	serverRepo := server.NewRepository(postgresSession)
	serverService := server.NewService(serverRepo)

	return &Application{
		PostgresSession: postgresSession,
		ServerService:   serverService,
	}, nil
}
