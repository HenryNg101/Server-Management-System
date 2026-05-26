package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	_ "github.com/HenryNg101/server-management-system/docs"
	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/feature/server"
	"github.com/HenryNg101/server-management-system/internal/feature/user"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/*
	TODO: Develop these APIs:

- POST   /api/v1/servers/import: Import servers from csv
- GET    /api/v1/servers/export: Export servers to csv
- POST   /api/v1/reports: Write report of servers statuses
*/
func SetupRouter(cfg *config.ApplicationConfig, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	user.RegisterRoutes(api, db)
	server.RegisterRoutes(api, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
