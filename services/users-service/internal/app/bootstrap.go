package app

import (
	"errors"

	"github.com/HenryNg101/users-service/internal/config"
	"github.com/HenryNg101/users-service/internal/platform/postgres"
	"github.com/HenryNg101/users-service/internal/user"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	UserService     user.Service
}

func NewApp() (*Application, error) {
	// Load config
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	userRepo := user.NewRepository(postgresSession)
	userService := user.NewService(userRepo)

	return &Application{
		PostgresSession: postgresSession,
		UserService:     userService,
	}, nil
}
