package auth

import "github.com/gin-gonic/gin"

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleRaw, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		userRole, ok := userRoleRaw.(string)
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}

		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
	}
}
