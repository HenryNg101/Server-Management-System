package data_transfer

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// Import servers from csv
// @Summary Import servers from csv file
// @Description User uploads an Excel file with servers info, and the server will try to import valid servers
// @Tags servers
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Input servers file (CSV)"
// @Success 200 {object} ImportServersResponse
// @Router /servers/import [post]
func (h *Handler) ImportServers(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"Input file error": "file is required"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{"Input file error": "cannot open file"})
		return
	}
	defer f.Close()

	result, err := h.service.CreateImportJob(c.Request.Context(), f, file.Size)
	if err != nil {
		c.JSON(500, gin.H{"Import servers error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// GET /jobs/:id
func (h *Handler) GetJob(c *gin.Context) {
	id := c.Param("id")

	job, err := h.service.GetJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
