package app

import (
	"github.com/HenryNg101/agents-service/internal/agent"
	"github.com/HenryNg101/agents-service/internal/middleware/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	agentH := agent.NewHandler(app.AgentService)

	// Agent-only routes
	agentGroup := rg.Group("/")
	agentGroup.Use(auth.AgentAuthMiddleware(app.AgentService, app.RedisSession))
	{
		agentGroup.POST("/metrics", agentH.IngestMetrics)
		agentGroup.POST("/rotate-key", agentH.RotateAPIKey)
	}

	// User-auth routes
	userGroup := rg.Group("/")
	userGroup.Use(auth.UserAuthMiddleware())
	{
		userGroup.POST("/register", agentH.RegisterAgent)
	}
}
