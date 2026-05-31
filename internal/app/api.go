package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_ "github.com/HenryNg101/server-management-system/docs"
	"github.com/HenryNg101/server-management-system/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO: Write more proper OpenAPI docs for the APIs
// TODO: Add JWT authentication to the system

// @title Server Management API
// @version 1.0
// @description API for managing servers and users, and reporting
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB, application *Application) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	RegisterRoutes(api, db, application)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
