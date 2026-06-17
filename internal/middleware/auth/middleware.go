package auth

import (
	"strings"

	"github.com/HenryNg101/server-management-system/internal/feature/agent"
	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing auth token"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token format"})
			return
		}

		claims, err := ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func AgentAuthMiddleware(db *gorm.DB, agentService agent.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing auth header"})
			return
		}

		// Expect: Bearer <API_KEY>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid auth format"})
			return
		}

		rawKey := parts[1]
		var agent model.Agent
		if err := agentService.AgentExist(c, rawKey, &agent); err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"})
			return
		}

		// attach server_id to context
		c.Set("server_id", agent.ServerID)

		c.Next()
	}
}
