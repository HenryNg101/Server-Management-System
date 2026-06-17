package agent

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// Help people log in, and generate tokens for them
// Metrics godoc
// @Summary Metrics
// @Description Ingest metrics from agent upload
// @Tags agents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body MetricMessage true "Server's metrics"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /agent/metrics [post]
func (h *Handler) IngestMetrics(c *gin.Context) {
	var msg MetricMessage

	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// TODO: override server_id
	serverID := c.GetUint("server_id")
	msg.ServerID = int(serverID)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
