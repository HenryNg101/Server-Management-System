package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	users := rg.Group("/servers")

	users.GET("/", handler.GetServers)
	users.POST("/", handler.CreateServer)
}
