package monitoring

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// @Summary Get statuses of servers and emailing
// @Description Get status report on the servers, and then send emails to people in the email list
// @Tags servers
// @Security BearerAuth
// @Produce json
// @Param request body SendReportRequest true "Report request"
// @Success 200 {object} Report
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /servers/report [post]
func (h *Handler) SendReports(c *gin.Context) {
	var req SendReportRequest
	var err error
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
		return
	}

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()
	if req.Start != nil {
		startTime, err = time.Parse("2006-01-02", *req.Start)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
			return
		}
	}
	if req.End != nil {
		endTime, err = time.Parse("2006-01-02", *req.End)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": err.Error()})
			return
		}
		endTime = endTime.Add(24 * time.Hour) // include the whole day
	}
	if req.TopN == nil {
		topN := 10
		req.TopN = &topN
	}
	if *req.TopN < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"Request parsing error": "Count must be greater than 0"})
		return
	}

	report, err := h.service.SendReports(
		startTime,
		endTime,
		*req.TopN,
		req.Emails,
		c.Request.Context(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}
