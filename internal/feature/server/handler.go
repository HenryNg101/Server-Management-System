package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) CreateServer(c *gin.Context) {
	var req CreateServerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.CreateServer(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// TODO: Add more advanced searches to this
func (h *Handler) GetServers(c *gin.Context) {
	servers, err := h.service.GetServers()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, servers)
}

func (h *Handler) GetServer(c *gin.Context) {
	var server model.Server

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	server, err = h.service.GetServer(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, server)
}

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
