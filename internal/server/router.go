package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/handler"
	"github.com/HenryNg101/server-management-system/internal/repository"
	"github.com/HenryNg101/server-management-system/internal/service"
)

func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Dependency chain
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	api := r.Group("/api/v1")

	users := api.Group("/users")
	{
		users.GET("/", userHandler.GetUsers)
		users.POST("/", userHandler.CreateUser)
	}

	return r
}
