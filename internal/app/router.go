package app

import (
	"github.com/HenryNg101/server-management-system/internal/feature/agent"
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
	agentHandler := agent.NewHandler(*application.AgentService)

	//
	// Agent's APIs
	agentGroup := rg.Group("/api/v1/agent")
	agentGroup.Use(internalAuth.AgentAuthMiddleware(*application.AgentService))
	{
		agentGroup.POST("/metrics", agentHandler.IngestMetrics)
	}

	//
	// Public API
	rg.POST("/login", authH.Login)
	rg.POST("/refresh", authH.Refresh)
	rg.POST("/logout", authH.Logout)

	//
	// Protected APIs
	protected := rg.Group("/")
	protected.Use(internalAuth.UserAuthMiddleware())

	protected.GET("/servers", serverHandler.GetServers)
	protected.GET("/servers/:id", serverHandler.GetServer)
	protected.GET("/users", userHandler.GetUsers)

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
