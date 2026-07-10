package app

import (
	"github.com/HenryNg101/monitoring-service/internal/monitoring"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	monitoringH := monitoring.NewHandler(app.MonitoringService)

	rg.POST("/monitoring/report", monitoringH.SendReports)
}
