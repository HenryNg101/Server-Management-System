package jobs

import (
	"log"
	"net/http"

	"github.com/HenryNg101/jobs-service/internal/model"
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
// @Success 200 {object} CreateImportJobResponse
// @Router /import-server [post]
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

// Get a job's status
// @Summary Get a job's status
// @Description After user uploads Excel file(s) with servers info for import, they can take the job's ID and check the status of processing in here
// @Param id path string true "Job ID"
// @Tags jobs
// @Security BearerAuth
// @Produce json
// @Success 200 {object} GetJobResponse
// @Failure 404 {object} map[string]string
// @Router /{id} [get]
func (h *Handler) GetJob(c *gin.Context) {
	id := c.Param("id")

	job, err := h.service.GetJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	ctx := c.Request.Context()

	resp := GetJobResponse{
		ID:               job.ID,
		Status:           string(job.Status),
		ProcessedRows:    job.ProcessedRows,
		SuccessRowsCount: job.SuccessRowsCount,
		FailedRowsCount:  job.FailedRowsCount,
		Error:            job.Error,
	}

	// Always allow download of input file
	if job.FilePath != "" {
		url, err := h.service.GenerateFileDownloadURL(ctx, job.FilePath)
		if err == nil {
			resp.InputFileURL = &url
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Only expose failures file when job is done
	if job.Status == model.JobStatusDone && job.ResultPath != nil {
		url, err := h.service.GenerateFileDownloadURL(ctx, *job.ResultPath)
		if err == nil {
			resp.FailuresFileURL = &url
		} else {
			log.Printf("Failures file does not exist yet: %v", err.Error())
		}
	}

	c.JSON(http.StatusOK, resp)
}
