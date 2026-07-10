package app

import (
	"github.com/HenryNg101/user-service/internal/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	userH := user.NewHandler(app.UserService)

	rg.GET("/", userH.GetUsers)
	rg.POST("/", userH.CreateUser)
}
