package app

import (
	"github.com/HenryNg101/user-service/internal/middleware/auth"
	"github.com/HenryNg101/user-service/internal/user"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	userH := user.NewHandler(app.UserService)

	protected := rg.Group("/")
	protected.Use(auth.UserAuthMiddleware())

	rg.GET("/", userH.GetUsers)

	admin := protected.Group("/")
	admin.Use(auth.RequireRoles("admin"))

	rg.POST("/", userH.CreateUser)
}
