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

	admin := protected.Group("/")
	admin.Use(auth.RequireRoles("admin"))

	admin.GET("/", userH.GetUsers)
	admin.POST("/", userH.CreateUser)
}
