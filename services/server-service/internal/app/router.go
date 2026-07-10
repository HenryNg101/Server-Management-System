package app

import (
	"github.com/HenryNg101/server-service/internal/middleware/auth"
	"github.com/HenryNg101/server-service/internal/server"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	serverH := server.NewHandler(app.ServerService)

	protected := rg.Group("/")
	protected.Use(auth.UserAuthMiddleware())

	protected.GET("/", serverH.GetServers)
	protected.GET("/:id", serverH.GetServer)

	admin := protected.Group("/")
	admin.Use(auth.RequireRoles("admin"))

	rg.POST("/", serverH.CreateServer)
	rg.PATCH("/:id", serverH.UpdateServer)
	rg.DELETE("/:id", serverH.DeleteServer)
	rg.GET("/export", serverH.ExportServers)
}
