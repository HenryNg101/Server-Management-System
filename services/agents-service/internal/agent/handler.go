package agent

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type UserContext struct {
	UserID uint
	Role   string
}

func GetUserContext(c *gin.Context) *UserContext {
	var userContext UserContext

	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return nil
	}
	userID, ok := userIDRaw.(uint)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid user context"})
		return nil
	}

	userRoleRaw, exists := c.Get("role")
	if !exists {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return nil
	}
	userRole, ok := userRoleRaw.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "invalid role context"})
		return nil
	}

	userContext.UserID = userID
	userContext.Role = userRole
	return &userContext
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

// RotateAPIKey godoc
// @Summary Rotate agent API key
// @Description Rotate API key for the authenticated agent. The old key will be invalidated immediately.
// @Tags agents
// @Produce json
// @Success 200 {object} map[string]string "New API key"
// @Failure 401 {object} map[string]string "Unauthorized - invalid or missing API key"
// @Failure 500 {object} map[string]string "Internal server error"
// @Security ApiKeyAuth
// @Router /rotate-key [post]
func (h *Handler) RotateAPIKey(c *gin.Context) {
	apiKey := c.GetHeader("X-Agent-API-Key")

	newKey, err := h.service.RotateAPIKey(c, apiKey)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"api_key": newKey,
	})
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

	userContext := GetUserContext(c)
	if userContext == nil {
		return
	}
	log.Println(userContext)

	resp, err := h.service.RegisterAgent(c, userContext.UserID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, resp)
}

// TODO: Agent un-registration
