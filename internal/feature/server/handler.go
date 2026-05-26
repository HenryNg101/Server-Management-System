package server

import (
	"encoding/csv"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/gin-gonic/gin"
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
// @Produce json
// @Success 200 {object} model.Server
// @Router /users [post]
func (h *Handler) CreateServer(c *gin.Context) {
	var req CreateServerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateServer(c.Request.Context(), req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// TODO: Add more advanced searches to this
// GetServers godoc
// @Summary Get servers
// @Description Retrieve list of suitable servers based on filters
// @Tags servers
// @Produce json
// @Success 200 {array} model.Server
// @Router /servers [get]
func (h *Handler) GetServers(c *gin.Context) {
	servers, err := h.service.GetServers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, servers)
}

// TODO: Add the logic to return 404 when no server with associated ID is found
// GetServer godoc
// @Summary Get a server
// @Description Retrieve a server based on ID
// @Tags servers
// @Produce json
// @Success 200 {object} model.Server
// @Router /servers/:id [get]
func (h *Handler) GetServer(c *gin.Context) {
	var server *model.Server

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	server, err = h.service.GetServer(c.Request.Context(), uint(id), server)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, server)
}

// TODO: Add the logic to return 404 when no server with associated ID is found
// UpdateServer godoc
// @Summary Update a server
// @Description Get a server based on ID, and then update field(s)
// @Tags servers
// @Produce json
// @Success 200 {object} model.Server
// @Router /servers/:id [patch]
func (h *Handler) UpdateServer(c *gin.Context) {
	var req UpdateServerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serverId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.UpdateServer(
		c.Request.Context(),
		uint(serverId),
		req,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteServer godoc
// @Summary Delete a server
// @Description Get a server based on ID, and delete that server
// @Tags servers
// @Produce json
// @Success 204
// @Router /servers/:id [delete]
func (h *Handler) DeleteServer(c *gin.Context) {
	serverId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.service.DeleteServer(c.Request.Context(), uint(serverId))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(404, gin.H{"error": "server not found"})
			return
		}

		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Status(204) // No Content
}

// Import servers from csv
// @Summary Import servers from csv file
// @Description User uploads an Excel file with servers info, and the server will try to import valid servers
// @Tags servers
// @Produce json
// @Success 200 {object} ImportServersResponse
// @Router /servers/import [post]
func (h *Handler) ImportServers(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "cannot open file"})
		return
	}
	defer f.Close()

	result, err := h.service.ImportServers(c.Request.Context(), f)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// Export servers to csv file
// @Summary Export servers info to csv file
// @Description User ask for all servers info, and the server will export servers info to CSV file
// @Tags servers
// @Produce text/csv
// @Success 200 {file} ImportServersResponse
// @Router /servers/export [get]
func (h *Handler) ExportServers(c *gin.Context) {
	servers, err := h.service.GetServers(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=servers.csv")

	writer := csv.NewWriter(c.Writer)

	// Write header
	header := []string{"id", "name", "status", "ipv4_address", "port", "protocol", "created_at", "last_updated_at"}
	if err := writer.Write(header); err != nil {
		c.JSON(500, gin.H{"error": "failed to write csv"})
		return
	}

	// Write rows
	for _, s := range servers {
		row := []string{
			strconv.FormatUint(uint64(s.ID), 10),
			s.Name,
			strconv.FormatBool(s.Status),
			s.IPv4Address,
			strconv.FormatUint(uint64(s.Port), 10),
			s.Protocol,
			s.CreatedAt.Format(time.RFC3339),
			s.LastUpdated.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			c.JSON(500, gin.H{"error": "failed to write csv row"})
			return
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}
