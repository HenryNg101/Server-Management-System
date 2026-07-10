package server

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/HenryNg101/server-service/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// CreateServer godoc
// @Summary Create a server
// @Description Create a new server, with specifications
// @Tags servers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateServerRequest true "Server payload"
// @Success 201 {object} CreateServerResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servers [post]
func (h *Handler) CreateServer(c *gin.Context) {
	var req CreateServerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
		return
	}

	createdServer, err := h.service.CreateServer(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Server creation error": err.Error()})
		return
	}
	createdAgent, apiKey, err := h.service.CreateAgent(c, createdServer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Agent creation error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, CreateServerResponse{
		CreatedServer: *createdServer,
		CreatedAgent:  *createdAgent,
		ApiKey:        apiKey,
	})
}

// @Summary Get servers
// @Description Retrieve servers with filtering, pagination, and sorting
// @Tags servers
// @Security BearerAuth
// @Produce json
// @Param status query bool false "Filter by status"
// @Param protocol query string false "Filter by protocol"
// @Param name query string false "Search by name"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10, max 100)"
// @Param sort_by query string false "Sort by field (id, name, created_at)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {array} GetServerResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servers [get]
func (h *Handler) GetServers(c *gin.Context) {
	query, err := ParseServersQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
	}
	paginatedResult, err := h.service.GetServers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Get servers error": err.Error()})
		return
	}
	result := make([]GetServerResponse, 0)
	for _, server := range paginatedResult.Servers {
		result = append(result, GetServerResponse{
			Name:        server.Name,
			Status:      server.Status,
			IPv4Address: server.IPv4Address,
			Port:        server.Port,
			Protocol:    server.Protocol,
			CreatedAt:   server.CreatedAt,
			LastUpdated: server.LastUpdated,
		})
	}
	c.JSON(http.StatusOK, result)
}

// GetServer godoc
// @Summary Get a server
// @Description Retrieve a server based on ID
// @Param id path int true "Server ID"
// @Tags servers
// @Security BearerAuth
// @Produce json
// @Success 200 {object} GetServerResponse
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servers/{id} [get]
func (h *Handler) GetServer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
		return
	}

	var server *model.Server
	server, err = h.service.GetServer(c.Request.Context(), uint(id), server)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"Get server error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Get server error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, &GetServerResponse{
		Name:        server.Name,
		Status:      server.Status,
		IPv4Address: server.IPv4Address,
		Port:        server.Port,
		Protocol:    server.Protocol,
		CreatedAt:   server.CreatedAt,
		LastUpdated: server.LastUpdated,
	})
}

// UpdateServer godoc
// @Summary Update a server
// @Description Get a server based on ID, and then update field(s)
// @Param id path int true "Server ID"
// @Accept json
// @Tags servers
// @Security BearerAuth
// @Produce json
// @Param request body UpdateServerRequest true "New server info"
// @Success 200 {object} model.Server
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /servers/{id} [patch]
func (h *Handler) UpdateServer(c *gin.Context) {
	var req UpdateServerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request body error": err.Error()})
		return
	}

	serverId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Parse server ID error": err.Error()})
		return
	}

	updated, err := h.service.UpdateServer(
		c.Request.Context(),
		uint(serverId),
		req,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"Update server error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Update server error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteServer godoc
// @Summary Delete a server
// @Description Get a server based on ID, and delete that server
// @Param id path int true "Server ID"
// @Tags servers
// @Security BearerAuth
// @Produce json
// @Success 204
// @Failure 500 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /servers/{id} [delete]
func (h *Handler) DeleteServer(c *gin.Context) {
	serverId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Parse server ID error": err.Error()})
		return
	}

	err = h.service.DeleteServer(c.Request.Context(), uint(serverId))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"Delete server error": "server not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"Delete server error": err.Error()})
		return
	}

	c.Status(204) // No Content
}

// Export servers to csv file
// @Summary Export servers info to csv file
// @Description User ask for all servers info with optional filtering, pagination, and sorting, and the server will export servers info to CSV file
// @Tags servers
// @Security BearerAuth
// @Produce text/csv
// @Param status query bool false "Filter by status"
// @Param protocol query string false "Filter by protocol"
// @Param name query string false "Search by name"
// @Param page query int false "Page number (default 1)"
// @Param page_size query int false "Page size (default 10, max 100)"
// @Param sort_by query string false "Sort by field (id, name, created_at)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {file} file "CSV file"
// @Header 200 {string} Content-Disposition "attachment; filename=servers.csv"
// @Router /servers/export [get]
func (h *Handler) ExportServers(c *gin.Context) {
	query, err := ParseServersQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
		return
	}

	servers, err := h.service.GetServers(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Get servers error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=servers.csv")

	writer := csv.NewWriter(c.Writer)

	// Write header
	header := []string{"name", "status", "ipv4_address", "port", "protocol", "created_at", "last_updated_at"}
	if err := writer.Write(header); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"File write error": err.Error()})
		return
	}

	// Write rows
	for _, s := range servers.Servers {
		row := []string{
			s.Name,
			strconv.FormatBool(s.Status),
			s.IPv4Address,
			strconv.FormatUint(uint64(s.Port), 10),
			s.Protocol,
			s.CreatedAt.Format(time.RFC3339),
			s.LastUpdated.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"File write error": err.Error()})
			return
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"File write error": err.Error()})
		return
	}
}
