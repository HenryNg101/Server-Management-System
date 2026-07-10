package app

import (
	"github.com/HenryNg101/jobs-service/internal/jobs"
	"github.com/HenryNg101/jobs-service/internal/middleware/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	jobH := jobs.NewHandler(app.JobsService)

	protected := rg.Group("/")
	protected.Use(auth.UserAuthMiddleware())

	admin := protected.Group("/")
	admin.Use(auth.RequireRoles("admin"))

	admin.POST("/import-server", jobH.ImportServers)
	admin.GET("/:id", jobH.GetJob)
}
