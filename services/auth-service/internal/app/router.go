package app

import (
	"github.com/HenryNg101/auth-service/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(rg *gin.RouterGroup, app *Application) {
	authH := auth.NewHandler(app.AuthService)

	rg.POST("/login", authH.Login)
	rg.POST("/refresh", authH.Refresh)
	rg.POST("/logout", authH.Logout)
}
