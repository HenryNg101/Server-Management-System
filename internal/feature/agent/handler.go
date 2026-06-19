package agent

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Ingest metrics from agents upload
// Metrics godoc
// @Summary Metrics
// @Description Ingest metrics from agent upload
// @Tags agents
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body []MetricMessage true "Server's metrics"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /agent/metrics [post]
func (h *Handler) IngestMetrics(c *gin.Context) {
	var msgs []MetricMessage

	if err := c.ShouldBindJSON(&msgs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	serverID := c.GetUint("server_id")

	for i := range msgs {
		msgs[i].ServerID = int(serverID)
	}

	// TODO: push msgs to Kafka
	fmt.Printf("Received %d metrics\n", len(msgs))

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
