package app

import (
	"github.com/HenryNg101/jobs-service/internal/jobs"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	jobH := jobs.NewHandler(app.JobsService)

	rg.POST("/import-server", jobH.ImportServers)
	rg.GET("/:id", jobH.GetJob)
}
