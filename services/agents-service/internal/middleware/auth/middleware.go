package auth

import (
	"github.com/HenryNg101/agent-service/internal/agent"
	"github.com/HenryNg101/agent-service/internal/model"
	"github.com/gin-gonic/gin"
)

func AgentAuthMiddleware(agentService agent.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Agent-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing auth header"})
			return
		}

		var agent model.Agent
		if err := agentService.AgentExist(c, apiKey, &agent); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"})
			return
		}

		// attach server_id to context
		c.Set("server_id", agent.ServerID)

		c.Next()
	}
}
