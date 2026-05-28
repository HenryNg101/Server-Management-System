package app

import (
	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	"github.com/HenryNg101/server-management-system/internal/platform/database"
	"gorm.io/gorm"
)

type Application struct {
	DB            *gorm.DB
	ServerService *server.Service
	UserService   *user.Service
}

// Create a new app with services to be used
// Why I do this ? So that, when an app is created, it can use all internal functionalities without having to use other ways of communication like gRPC or HTTP
// If it's API app, you can add handlers, add API group, etc. to it. If it's a worker, it can still get access to DB query, services, etc. without it being separated entirely
// It might be tightly coupled, but it's good for now
func NewApp() (*Application, error) {
	postgresConfig := config.LoadPostgres()
	postgresSession := database.NewPostgresSession(postgresConfig)

	serverRepo := server.NewRepository(postgresSession)
	serverService := server.NewService(serverRepo)

	userRepo := user.NewRepository(postgresSession)
	userService := user.NewService(userRepo)

	return &Application{
		DB:            postgresSession,
		ServerService: &serverService,
		UserService:   &userService,
	}, nil
}
