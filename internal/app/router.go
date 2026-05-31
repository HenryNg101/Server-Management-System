package app

import (
	authHandler "github.com/HenryNg101/server-management-system/internal/feature/auth"
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	internalAuth "github.com/HenryNg101/server-management-system/internal/middleware/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, application *Application) {
	// Register auth handler for authentication
	authH := authHandler.NewHandler(*application.AuthService)

	// Register handlers for servers and users
	serverHandler := server.NewHandler(*application.ServerService)
	userHandler := user.NewHandler(*application.UserService)

	//
	// Public API
	rg.POST("/login", authH.Login)

	//
	// Protected APIs
	protected := rg.Group("/")
	protected.Use(internalAuth.AuthMiddleware())

	protected.GET("/servers", serverHandler.GetServers)
	// protected.POST("/servers", serverHandler.CreateServer)
	protected.GET("/servers/:id", serverHandler.GetServer)
	// protected.PATCH("/servers/:id", serverHandler.UpdateServer)
	// protected.DELETE("/servers/:id", serverHandler.DeleteServer)
	// protected.POST("/servers/import", serverHandler.ImportServers)
	// protected.GET("/servers/export", serverHandler.ExportServers)
	// protected.POST("/servers/report", serverHandler.SendReports)

	protected.GET("/users", userHandler.GetUsers)
	// protected.POST("/users", userHandler.CreateUser)

	//
	// Admin only APIs
	admin := protected.Group("/")
	admin.Use(internalAuth.RequireRoles("admin"))

	admin.POST("/servers", serverHandler.CreateServer)
	admin.PATCH("/servers/:id", serverHandler.UpdateServer)
	admin.DELETE("/servers/:id", serverHandler.DeleteServer)
	admin.POST("/servers/import", serverHandler.ImportServers)
	admin.POST("/users", userHandler.CreateUser)
	admin.GET("/servers/export", serverHandler.ExportServers)
	admin.POST("/servers/report", serverHandler.SendReports)
}
