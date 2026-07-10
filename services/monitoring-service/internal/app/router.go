package app

import (
	"github.com/HenryNg101/monitoring-service/internal/middleware/auth"
	"github.com/HenryNg101/monitoring-service/internal/monitoring"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	monitoringH := monitoring.NewHandler(app.MonitoringService)

	protected := rg.Group("/")
	protected.Use(auth.UserAuthMiddleware())

	admin := protected.Group("/")
	admin.Use(auth.RequireRoles("admin"))

	admin.POST("/report", monitoringH.SendReports)
}
