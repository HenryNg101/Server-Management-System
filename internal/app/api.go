package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	docs "github.com/HenryNg101/server-management-system/docs"
	"github.com/HenryNg101/server-management-system/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO: Write more proper OpenAPI docs for the APIs

// @title Server Management API
// @version 1.0
// @description API for managing servers and users, and reporting
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your JWT token.
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-Agent-API-Key
// @description API Key for external automated systems.
func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB, application *Application) *gin.Engine {
	r := gin.Default()

	// Dynamic Swagger config
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port
	// docs.SwaggerInfo.BasePath = "/api/v1"
	// docs.SwaggerInfo.Schemes = []string{"http"}

	api := r.Group("/api/v1")
	RegisterRoutes(api, db, application)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
