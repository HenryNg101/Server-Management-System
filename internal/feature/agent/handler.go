package agent

import (
	"github.com/HenryNg101/server-management-system/internal/platform/kafka"
	"github.com/gin-gonic/gin"
)

func IngestMetrics(producer *kafka.Producer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg MetricMessage

		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(400, gin.H{"error": "invalid payload"})
			return
		}

		// 🔥 CRITICAL: override server_id
		serverID := c.GetUint("server_id")
		msg.ServerID = int(serverID)

		// send to Kafka
		producer.Send(msg)

		c.JSON(200, gin.H{"status": "ok"})
	}
}
