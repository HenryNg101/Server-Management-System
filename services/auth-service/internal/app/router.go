package app

import (
	"github.com/HenryNg101/auth-service/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	authH := auth.NewHandler(app.AuthService)

	rg.POST("/login", authH.Login)
	rg.POST("/refresh", authH.Refresh)
	rg.POST("/logout", authH.Logout)
}

// func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, application *Application) {
// 	// Register auth handler for authentication
// 	authH := authHandler.NewHandler(application.AuthService)

// 	// Register handlers for servers and users
// 	serverHandler := server.NewHandler(application.ServerService)
// 	userHandler := user.NewHandler(application.UserService)
// 	agentHandler := agent.NewHandler(application.AgentService)
// 	monitoringHandler := monitoring.NewHandler(application.MonitoringService)
// 	dataTransferHandler := data_transfer.NewHandler(application.DataTransferService)

// 	//
// 	// Agent's APIs
// 	agentGroup := rg.Group("/agent")
// 	agentGroup.Use(internalAuth.AgentAuthMiddleware(application.AgentService))
// 	{
// 		agentGroup.POST("/metrics", agentHandler.IngestMetrics)
// 	}

// 	//
// 	// Public API
// 	rg.POST("/login", authH.Login)
// 	rg.POST("/refresh", authH.Refresh)
// 	rg.POST("/logout", authH.Logout)

// 	//
// 	// Protected APIs
// 	protected := rg.Group("/")
// 	protected.Use(internalAuth.UserAuthMiddleware())

// 	protected.GET("/servers", serverHandler.GetServers)
// 	protected.GET("/servers/:id", serverHandler.GetServer)
// 	protected.GET("/users", userHandler.GetUsers)

// 	//
// 	// Admin only APIs
// 	admin := protected.Group("/")
// 	admin.Use(internalAuth.RequireRoles("admin"))

// 	admin.POST("/servers", serverHandler.CreateServer)
// 	admin.PATCH("/servers/:id", serverHandler.UpdateServer)
// 	admin.DELETE("/servers/:id", serverHandler.DeleteServer)
// 	admin.POST("/servers/import", dataTransferHandler.ImportServers)
// 	admin.GET("/jobs/:id", dataTransferHandler.GetJob)
// 	admin.POST("/users", userHandler.CreateUser)
// 	admin.GET("/servers/export", serverHandler.ExportServers)
// 	admin.POST("/servers/report", monitoringHandler.SendReports)
// }
