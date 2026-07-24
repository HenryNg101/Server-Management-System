package app

import (
	"github.com/HenryNg101/servers-service/internal/middleware/auth"
	"github.com/HenryNg101/servers-service/internal/server"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	serverH := server.NewHandler(app.ServerService)

	protected := rg.Group("/")
	protected.Use(auth.UserAuthMiddleware())

	protected.GET("/", serverH.GetServers)
	protected.GET("/:id", serverH.GetServer)
	protected.POST("/", serverH.CreateServer)
	protected.PATCH("/:id", serverH.UpdateServer)
	protected.DELETE("/:id", serverH.DeleteServer)
	protected.GET("/export", serverH.ExportServers)
}
