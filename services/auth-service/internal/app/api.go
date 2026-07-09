package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/HenryNg101/auth-service/docs"
	"github.com/HenryNg101/auth-service/internal/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service API
// @version 1.0
// @description API for authenticate/authorize users
// @host localhost
// @BasePath /auth
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your JWT token.
func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB, application *Application) *gin.Engine {
	r := gin.Default()

	// Dynamic Swagger config
	docs.SwaggerInfo.Host = cfg.Host + ":" + cfg.Port
	// docs.SwaggerInfo.BasePath = "/api/v1"
	// docs.SwaggerInfo.Schemes = []string{"http"}

	api := r.Group("/")
	RegisterRoutes(api, application)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
