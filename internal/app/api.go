package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
)

func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Dependency chain
	// userRepo := user.NewRepository(db)
	// userService := user.NewService(userRepo)
	// userHandler := user.NewHandler(userService)

	api := r.Group("/api/v1")

	// users := api.Group("/users")
	// {
	// 	users.GET("/", userHandler.GetUsers)
	// 	users.POST("/", userHandler.CreateUser)
	// }

	user.RegisterRoutes(api, db)

	return r
}
