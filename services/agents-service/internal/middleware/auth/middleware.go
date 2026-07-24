package auth

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/HenryNg101/agents-service/internal/agent"
	"github.com/HenryNg101/agents-service/internal/shared/auth"
	internalAuth "github.com/HenryNg101/agents-service/internal/shared/auth"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type LocalCache struct {
	data map[string]uint
	mu   sync.RWMutex
}

// TODO: Fix this
func AgentAuthMiddleware(agentService agent.Service, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-Agent-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing auth header"})
			return
		}

		hashedKey := auth.HashAPIKey(apiKey)
		redisKey := "agent:auth:" + hashedKey

		// 1. Try Redis first
		serverIDStr, err := redisClient.Get(c, redisKey).Result()
		if err == nil {
			serverID, _ := strconv.Atoi(serverIDStr)
			c.Set("server_id", uint(serverID))
			c.Next()
			return
		}

		// 2. If Redis miss → fallback to find in DB
		if err != redis.Nil {
			log.Println("[WARN] redis error:", err)
		}

		agentModel, err := agentService.FindByHashedKey(c, hashedKey)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid api key"})
			return
		}

		// 3. Backfill Redis
		if err := redisClient.Set(c, redisKey, agentModel.ServerID, 0).Err(); err != nil {
			log.Println("[WARN] redis backfill failed:", err)
		}

		// attach server_id to context
		c.Set("server_id", agentModel.ServerID)
		c.Next()
	}
}

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

		claims, err := internalAuth.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		// Refresh token is not for this, so it's not allowed here
		if claims.Type != "access" {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token type"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
