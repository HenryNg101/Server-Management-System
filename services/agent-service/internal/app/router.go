package app

import (
	"github.com/HenryNg101/agent-service/internal/agent"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	agentH := agent.NewHandler(app.AgentService)

	rg.POST("/metrics", agentH.IngestMetrics)
}
