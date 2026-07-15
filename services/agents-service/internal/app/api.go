package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HenryNg101/agent-service/docs"
	"github.com/HenryNg101/agent-service/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Agent Service API
// @version 1.0
// @description APIs for agents to push data to system
// @host localhost
// @BasePath /agents
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Agent-API-Key
// @description API Key for external automated systems.
func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB, application *Application) *gin.Engine {
	r := gin.Default()

	// Dynamic Swagger config
	docs.SwaggerInfo.Host = cfg.Host
	// docs.SwaggerInfo.BasePath = "/api/v1"
	// docs.SwaggerInfo.Schemes = []string{"http"}

	api := r.Group("/")
	RegisterRoutes(api, application)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
