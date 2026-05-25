package main

import (
	"log"

	"github.com/HenryNg101/server-management-system/internal/app"
	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/HenryNg101/server-management-system/internal/platform/database"
)

func main() {
	postgresConfig := config.LoadPostgres()
	appConfig := config.LoadAppConfig()

	postgresSession := database.NewPostgresSession(postgresConfig)

	router := app.SetupRouter(appConfig, postgresSession)

	if err := router.Run(":" + appConfig.Port); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
