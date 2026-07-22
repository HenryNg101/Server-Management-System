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
// @Router /metrics [post]
func (h *Handler) IngestMetrics(c *gin.Context) {
	var msgs []MetricMessage

	if err := c.ShouldBindJSON(&msgs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Limit the request body size to prevent abuse
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1<<20) // 1MB
	if len(msgs) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many metrics"})
		return
	}

	serverID := c.GetUint("server_id")

	for i := range msgs {
		msgs[i].ServerID = int(serverID)
	}

	if err := h.service.PushMetrics(c, msgs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to ingest metrics"})
		return
	}

	fmt.Printf("Received %d metrics\n", len(msgs))

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Register a new agent to the system
// RegisterAgent godoc
// @Summary Register a new agent
// @Description Register an agent instance to an existing server owned by the authenticated user
// @Tags agents
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body RegisterAgentRequest true "Agent registration payload"
// @Success 201 {object} RegisterAgentResponse
// @Failure 400 {object} map[string]string "Bad request / validation error"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /register [post]
func (h *Handler) RegisterAgent(c *gin.Context) {
	var req RegisterAgentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := userIDRaw.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user context"})
		return
	}

	resp, err := h.service.RegisterAgent(c, userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, resp)
}

// TODO: Agent un-registration
