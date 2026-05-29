package app

import (
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, application *Application) {
	// Register handlers for servers
	serverHandler := server.NewHandler(*application.ServerService)

	servers := rg.Group("/servers")

	servers.GET("/", serverHandler.GetServers)
	servers.POST("/", serverHandler.CreateServer)
	servers.GET("/:id", serverHandler.GetServer)
	servers.PATCH("/:id", serverHandler.UpdateServer)
	servers.DELETE("/:id", serverHandler.DeleteServer)
	servers.POST("/import", serverHandler.ImportServers)
	servers.GET("/export", serverHandler.ExportServers)
	servers.POST("/report", serverHandler.SendReports)

	// Register handlers for users
	userHandler := user.NewHandler(*application.UserService)

	users := rg.Group("/users")

	users.GET("/", userHandler.GetUsers)
	users.POST("/", userHandler.CreateUser)
}
