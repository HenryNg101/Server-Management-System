package app

import (
	"github.com/HenryNg101/server-service/internal/server"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	serverH := server.NewHandler(app.ServerService)

	rg.GET("/", serverH.GetServers)
	rg.GET("/:id", serverH.GetServer)

	rg.POST("/", serverH.CreateServer)
	rg.PATCH("/:id", serverH.UpdateServer)
	rg.DELETE("/:id", serverH.DeleteServer)
	rg.GET("/export", serverH.ExportServers)
}
